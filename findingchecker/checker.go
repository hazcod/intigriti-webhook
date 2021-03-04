package findingchecker

import (
	"github.com/hazcod/go-intigriti"
	"github.com/hazcod/intigriti-webhook/config"
	"github.com/hazcod/intigriti-webhook/webhook"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"net/url"
	"time"
)

func schedule(what func(), delay time.Duration, stopChan chan bool) {
	logrus.WithField("check_interval_minutes", delay.Minutes()).Debug("checks scheduled")
	ticker := time.NewTicker(delay)
	go func() {
		for {
			select {
			case <-stopChan:
				logrus.Debug("stopping checker")
				return
			case <-ticker.C:
				what()
			}
		}
	}()
}

func findingExists(config config.Config, finding intigriti.Submission) bool {
	for _, fID := range config.FindingIDs {
		if fID == finding.ID {
			return true
		}
	}

	return false
}

func savetoConfig(config config.Config, findings []intigriti.Submission) error {
	bytes, err := yaml.Marshal(config)
	if err != nil {
		return errors.Wrap(err, "could not marshal config yaml")
	}

	return ioutil.WriteFile(config.ConfigPath, bytes, 0644)
}

func checkForNew(config config.Config, slckEndpoint webhook.Endpoint, intiEndpoint intigriti.Endpoint) (func(), error) {
	return func(){
		logrus.Debug("checking for new findings")
		findings, err := intiEndpoint.GetSubmissions()
		if err != nil {
			logrus.WithError(err).Error("could not fetch from intigriti")
			return
		}

		if len(findings) == 0 {
			logrus.Debug("no findings found")
			return
		}

		for _, finding := range findings {
			fLogger := logrus.WithField("finding_id", finding.ID).WithField("finding_state", finding.State)

			if config.IncludeNonReady && !finding.IsReady() {
				fLogger.Debug("skipping non-ready finding")
				continue
			}

			fLogger.Debug("looking if finding exists")
			if findingExists(config, finding) {
				fLogger.Debug("finding already sent to slack, skipping")
				continue
			}

			fLogger.Debug("new finding, sending off to webhook")
			if err := slckEndpoint.Send(finding); err != nil {
				logrus.WithError(err).Error("could not send to webhook")
			} else {
				config.FindingIDs = append(config.FindingIDs, finding.ID)
			}
		}

		logrus.WithField("findings_size", len(findings)).Debug("saving findings to our config")
		if err := savetoConfig(config, findings); err != nil {
			logrus.WithError(err).Error("could not add finding ID to config")
		}
	}, nil
}

func RunChecker(config config.Config, clientVersion string) error {
	webhookURL, err := url.Parse(config.WebhookURL)
	if err != nil {
		return errors.Wrap(err, "invalid slack url")
	}

	slackEndpoint := webhook.NewEndpoint(webhookURL, config.HTTPMethod, config.WebhookHeaders, config.Format, clientVersion)
	intigritiEndpoint := intigriti.New(config.IntigritiClientID, config.IntigritiClientSecret)

	checkFunc, err := checkForNew(config, slackEndpoint, intigritiEndpoint)
	if err != nil {
		return errors.Wrap(err, "could not initialize checker")
	}

	// should we ever want to stop it
	stopChan := make(chan bool, 1)

	// recurring runs
	schedule(checkFunc, time.Minute * time.Duration(config.CheckInterval), stopChan)

	logrus.Info("checker is is now running")

	// trigger first run immediately
	checkFunc()

	return nil
}
