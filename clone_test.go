package main

import "testing"

func TestCloneProject(t *testing.T) {
	type args struct {
		tenant   string
		repoName string
		repoUrl  string
		id       string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Default",
			args: args{
				tenant:   "theykk-bunker",
				repoName: "theykk-bunker/tester",
				repoUrl:  "https://github.com/theykk-bunker/tester.git",
				id:       "1",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CloneProject(tt.args.tenant, tt.args.repoName, tt.args.repoUrl, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("CloneProject() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
