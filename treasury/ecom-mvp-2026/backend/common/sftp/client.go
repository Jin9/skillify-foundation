package sftp

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// Client interface defines the available operations over SFTP.
type Client interface {
	Upload(ctx context.Context, remotePath string, content io.Reader) error
	Download(ctx context.Context, remotePath string, out io.Writer) error
	ListDirectory(ctx context.Context, remotePath string) ([]os.FileInfo, error)
	Delete(ctx context.Context, remotePath string) error
	Close() error
}

type client struct {
	sshClient  *ssh.Client
	sftpClient *sftp.Client
}

// NewClient constructs a new SSH connection and opens an SFTP session.
func NewClient(cfg Config) (Client, error) {
	if cfg.Timeout == 0 {
		cfg.Timeout = 10 * time.Second
	}

	authMethods := make([]ssh.AuthMethod, 0)
	if len(cfg.PrivateKeyAuth) > 0 {
		signer, err := ssh.ParsePrivateKey(cfg.PrivateKeyAuth)
		if err != nil {
			return nil, fmt.Errorf("failed to parse private key: %w", err)
		}
		authMethods = append(authMethods, ssh.PublicKeys(signer))
	} else if cfg.Password != "" {
		authMethods = append(authMethods, ssh.Password(cfg.Password))
	} else {
		return nil, errors.New("no password or private key provided for sftp")
	}

	sshConfig := &ssh.ClientConfig{
		User:            cfg.Username,
		Auth:            authMethods,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // Use cautiously in prod
		Timeout:         cfg.Timeout,
	}

	address := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	sshConn, err := ssh.Dial("tcp", address, sshConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to dial ssh: %w", err)
	}

	sftpConn, err := sftp.NewClient(sshConn)
	if err != nil {
		sshConn.Close()
		return nil, fmt.Errorf("failed to open sftp session: %w", err)
	}

	return &client{
		sshClient:  sshConn,
		sftpClient: sftpConn,
	}, nil
}

func (c *client) Close() error {
	var errs []error
	if err := c.sftpClient.Close(); err != nil {
		errs = append(errs, err)
	}
	if err := c.sshClient.Close(); err != nil {
		errs = append(errs, err)
	}
	if len(errs) > 0 {
		return fmt.Errorf("failed to close: %v", errs)
	}
	return nil
}

// Upload handles streaming content to the remote path.
func (c *client) Upload(ctx context.Context, remotePath string, content io.Reader) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	file, err := c.sftpClient.Create(remotePath)
	if err != nil {
		return fmt.Errorf("failed to create remote file: %w", err)
	}
	defer file.Close()

	done := make(chan error, 1)
	go func() {
		_, err := io.Copy(file, content)
		done <- err
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-done:
		return err
	}
}

// Download handles streaming content from the remote path to the out writer.
func (c *client) Download(ctx context.Context, remotePath string, out io.Writer) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	file, err := c.sftpClient.Open(remotePath)
	if err != nil {
		return fmt.Errorf("failed to open remote file: %w", err)
	}
	defer file.Close()

	done := make(chan error, 1)
	go func() {
		_, err := io.Copy(out, file)
		done <- err
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-done:
		return err
	}
}

// ListDirectory lists contents of a directory.
func (c *client) ListDirectory(ctx context.Context, remotePath string) ([]os.FileInfo, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	type listResult struct {
		files []os.FileInfo
		err   error
	}
	done := make(chan listResult, 1)
	go func() {
		files, err := c.sftpClient.ReadDir(remotePath)
		done <- listResult{files, err}
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case res := <-done:
		return res.files, res.err
	}
}

// Delete removes a remote file.
func (c *client) Delete(ctx context.Context, remotePath string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	done := make(chan error, 1)
	go func() {
		done <- c.sftpClient.Remove(remotePath)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-done:
		return err
	}
}
