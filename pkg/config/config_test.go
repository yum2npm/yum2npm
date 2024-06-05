package config

import (
	"reflect"
	"testing"
	"time"
)

func TestLoadConfigs(t *testing.T) {
	testCases := []struct {
		name           string
		path           string
		expectedResult Config
		expectedError  string
	}{
		{
			name: "Load a complete and valid config file",
			path: "testdata/complete.yaml",
			expectedResult: Config{
				HTTP: HTTP{
					ListenAddress: "0.0.0.0:8443",
					CertFile:      "/etc/yum2npm/fullchain.pem",
					KeyFile:       "/etc/yum2npm/privkey.pem",
				},
				Repos: []Repo{
					{
						Name: "rocky-9-appstream-x86_64",
						Url:  "https://dl.rockylinux.org/pub/rocky/9/AppStream/x86_64/os",
					},
					{
						Name: "rocky-9-baseos-x86_64",
						Url:  "https://dl.rockylinux.org/pub/rocky/9/BaseOS/x86_64/os",
					},
					{
						Name: "rocky-9-extras-x86_64",
						Url:  "https://dl.rockylinux.org/pub/rocky/9/extras/x86_64/os",
					},
					{
						Name: "rocky-9-crb-x86_64",
						Url:  "https://dl.rockylinux.org/pub/rocky/9/CRB/x86_64/os/",
					},
				},
				RefreshInterval: time.Hour,
			},
		},
		{
			name:           "Load an empty config file",
			path:           "testdata/empty.yaml",
			expectedResult: Config{},
			expectedError:  "no repos are configured",
		},
		{
			name:           "Load a config file with an invalid repo name",
			path:           "testdata/invalid_repo_name.yaml",
			expectedResult: Config{},
			expectedError:  "repo 0 does not match regex \"^[a-zA-Z0-9-_]+$\" (Rocky 9 AppStream (x86_64))",
		},
		{
			name:           "Load a config file with an invalid repo url",
			path:           "testdata/invalid_repo_url.yaml",
			expectedResult: Config{},
			expectedError:  "repo 0 (rocky-9-appstream-x86_64) has an invalid URL",
		},
		{
			name:           "Load a config file with a tls key but no certificate",
			path:           "testdata/no_cert.yaml",
			expectedResult: Config{},
			expectedError:  "http cert file must be specified if key file is",
		},
		{
			name:           "Load a config file with a tls cert but no key",
			path:           "testdata/no_key.yaml",
			expectedResult: Config{},
			expectedError:  "http key file must be specified if cert file is",
		},
		{
			name:           "Load a config file with a repo without an url",
			path:           "testdata/no_repo_url.yaml",
			expectedResult: Config{},
			expectedError:  "repo 0 (rocky-9-appstream-x86_64) does not have an URL",
		},
		{
			name: "Load a config file that only contains repos",
			path: "testdata/only_repos.yaml",
			expectedResult: Config{
				HTTP: HTTP{
					ListenAddress: ":8080",
				},
				Repos: []Repo{
					{
						Name: "rocky-9-appstream-x86_64",
						Url:  "https://dl.rockylinux.org/pub/rocky/9/AppStream/x86_64/os",
					},
					{
						Name: "rocky-9-baseos-x86_64",
						Url:  "https://dl.rockylinux.org/pub/rocky/9/BaseOS/x86_64/os",
					},
					{
						Name: "rocky-9-extras-x86_64",
						Url:  "https://dl.rockylinux.org/pub/rocky/9/extras/x86_64/os",
					},
					{
						Name: "rocky-9-crb-x86_64",
						Url:  "https://dl.rockylinux.org/pub/rocky/9/CRB/x86_64/os/",
					},
				},
				RefreshInterval: time.Hour,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c, err := Init(tc.path)
			if len(tc.expectedError) > 0 && err == nil {
				t.Errorf("expected an error but got none")
			}
			if err != nil && err.Error() != tc.expectedError {
				t.Errorf("expected error \"%s\" but got \"%v\"", tc.expectedError, err)
			}
			if !reflect.DeepEqual(tc.expectedResult, c) {
				t.Errorf("expected %v but got %v", tc.expectedResult, c)
			}
		})
	}
}
