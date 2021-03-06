package service

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
	"time"
)

type Service struct {
	URLTemplate string
	CmdTemplate []string
	url         string
	command     *exec.Cmd
}

type addressInfo struct {
	Address string
	Host    string
	Port    string
}

func (s *Service) URL() (string, error) {
	if s.command == nil {
		return "", errors.New("not running")
	}

	return s.url, nil
}

func (s *Service) Start() error {
	if s.command != nil {
		return errors.New("already running")
	}

	address, err := freeAddress()
	if err != nil {
		return fmt.Errorf("failed to locate a free port: %s", err)
	}

	if s.url, err = buildURL(s.URLTemplate, address); err != nil {
		return fmt.Errorf("failed to parse URL: %s", err)
	}

	command, err := buildCommand(s.CmdTemplate, address)
	if err != nil {
		return fmt.Errorf("failed to parse command: %s", err)
	}

	if err := command.Start(); err != nil {
		return fmt.Errorf("failed to run command: %s", err)
	}

	s.command = command

	return nil
}

func (s *Service) Stop() error {
	if s.command == nil {
		return errors.New("already stopped")
	}

	var err error
	if runtime.GOOS == "windows" {
		err = s.command.Process.Kill()
	} else {
		err = s.command.Process.Signal(syscall.SIGTERM)
	}
	if err != nil {
		return fmt.Errorf("failed to stop command: %s", err)
	}

	s.command.Wait()
	s.command = nil

	return nil
}

func freeAddress() (addressInfo, error) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return addressInfo{}, err
	}
	defer listener.Close()

	address := listener.Addr().String()
	addressParts := strings.SplitN(address, ":", 2)
	return addressInfo{address, addressParts[0], addressParts[1]}, nil
}

func (s *Service) WaitForBoot(timeout time.Duration) error {
	timeoutChan := time.After(timeout)
	failedChan := make(chan struct{}, 1)
	startedChan := make(chan struct{})

	go func() {
		up := s.checkStatus()
		for !up {
			select {
			case <-failedChan:
				return
			default:
				time.Sleep(500 * time.Millisecond)
				up = s.checkStatus()
			}
		}
		startedChan <- struct{}{}
	}()

	select {
	case <-timeoutChan:
		failedChan <- struct{}{}
		return errors.New("failed to start before timeout")
	case <-startedChan:
		return nil
	}
}

func (s *Service) checkStatus() bool {
	client := &http.Client{}
	request, _ := http.NewRequest("GET", fmt.Sprintf("%s/status", s.url), nil)
	response, err := client.Do(request)
	if err == nil && response.StatusCode == 200 {
		return true
	}
	return false
}
