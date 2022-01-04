package runner

type Runner interface {
	Format(repoName string) error
}
