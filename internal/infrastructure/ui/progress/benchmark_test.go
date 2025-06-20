package progress

import (
	"fmt"
	"testing"
	"time"

	"github.com/qskkk/git-fleet/internal/domain/entities"
)

// BenchmarkProgressBarRender benchmarks the rendering performance
func BenchmarkProgressBarRender(b *testing.B) {
	repositories := make([]string, 100)
	for i := 0; i < 100; i++ {
		repositories[i] = fmt.Sprintf("repo%d", i)
	}

	pb := NewProgressBar(repositories, "git status")

	// Pre-populate with some results
	for i := 0; i < 50; i++ {
		result := entities.NewExecutionResult(repositories[i], "git status")
		if i%3 == 0 {
			result.MarkAsFailed("error", 1, "test failure")
		} else {
			result.MarkAsSuccess("success", 0)
		}
		pb.UpdateProgress(result)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		output := pb.Render()
		_ = len(output) // Use output to prevent optimization
	}
}

// BenchmarkProgressBarRenderCompleted benchmarks rendering when completed
func BenchmarkProgressBarRenderCompleted(b *testing.B) {
	repositories := make([]string, 100)
	for i := 0; i < 100; i++ {
		repositories[i] = fmt.Sprintf("repo%d", i)
	}

	pb := NewProgressBar(repositories, "git status")

	// Complete all repositories
	for i := 0; i < 100; i++ {
		result := entities.NewExecutionResult(repositories[i], "git status")
		if i%4 == 0 {
			result.MarkAsFailed("error", 1, "test failure")
		} else {
			result.MarkAsSuccess("success", 0)
		}
		pb.UpdateProgress(result)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		output := pb.Render()
		_ = len(output) // Use output to prevent optimization
	}
}

// BenchmarkProgressBarUpdateProgress benchmarks updating progress
func BenchmarkProgressBarUpdateProgress(b *testing.B) {
	repositories := make([]string, b.N)
	for i := 0; i < b.N; i++ {
		repositories[i] = fmt.Sprintf("repo%d", i)
	}

	pb := NewProgressBar(repositories, "git status")

	results := make([]*entities.ExecutionResult, b.N)
	for i := 0; i < b.N; i++ {
		result := entities.NewExecutionResult(repositories[i], "git status")
		result.MarkAsSuccess("success", 0)
		results[i] = result
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pb.UpdateProgress(results[i])
	}
}

// BenchmarkProgressServiceConcurrentUpdates benchmarks concurrent updates
func BenchmarkProgressServiceConcurrentUpdates(b *testing.B) {
	service := &ProgressService{enabled: true}

	repositories := make([]string, 100)
	for i := 0; i < 100; i++ {
		repositories[i] = fmt.Sprintf("repo%d", i)
	}

	service.StartProgress(repositories, "git status")

	results := make([]*entities.ExecutionResult, 100)
	for i := 0; i < 100; i++ {
		result := entities.NewExecutionResult(repositories[i], "git status")
		result.MarkAsSuccess("success", 0)
		results[i] = result
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			service.UpdateProgress(results[i%100])
			i++
		}
	})
}

// BenchmarkNewProgressBar benchmarks creating new progress bars
func BenchmarkNewProgressBar(b *testing.B) {
	repositories := make([]string, 100)
	for i := 0; i < 100; i++ {
		repositories[i] = fmt.Sprintf("repo%d", i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewProgressBar(repositories, "git status")
	}
}

// BenchmarkProgressBarGetPercentage benchmarks percentage calculation
func BenchmarkProgressBarGetPercentage(b *testing.B) {
	repositories := make([]string, 1000)
	for i := 0; i < 1000; i++ {
		repositories[i] = fmt.Sprintf("repo%d", i)
	}

	pb := NewProgressBar(repositories, "git status")
	pb.completed = 500 // 50% completed

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = pb.GetPercentage()
	}
}

// BenchmarkProgressBarMarkRepositoryAsStarting benchmarks marking repositories as starting
func BenchmarkProgressBarMarkRepositoryAsStarting(b *testing.B) {
	repositories := make([]string, b.N)
	for i := 0; i < b.N; i++ {
		repositories[i] = fmt.Sprintf("repo%d", i)
	}

	pb := NewProgressBar(repositories, "git status")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pb.MarkRepositoryAsStarting(repositories[i])
	}
}

// BenchmarkProgressServiceWithDisabled benchmarks disabled progress service
func BenchmarkProgressServiceWithDisabled(b *testing.B) {
	service := &ProgressService{enabled: false}

	repositories := make([]string, 100)
	for i := 0; i < 100; i++ {
		repositories[i] = fmt.Sprintf("repo%d", i)
	}

	result := entities.NewExecutionResult("repo1", "git status")
	result.MarkAsSuccess("success", 0)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.StartProgress(repositories, "git status")
		service.MarkRepositoryAsStarting("repo1")
		service.UpdateProgress(result)
		service.FinishProgress()
	}
}

// BenchmarkStringBuilderVsFormatter benchmarks string building approaches
func BenchmarkStringBuilderVsFormatter(b *testing.B) {
	repositories := make([]string, 50)
	for i := 0; i < 50; i++ {
		repositories[i] = fmt.Sprintf("repo%d", i)
	}

	pb := NewProgressBar(repositories, "git status")

	// Add some completed results
	for i := 0; i < 25; i++ {
		result := entities.NewExecutionResult(repositories[i], "git status")
		result.MarkAsSuccess("success", 0)
		result.Duration = time.Duration(i) * time.Millisecond
		pb.UpdateProgress(result)
	}

	b.Run("current_implementation", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = pb.Render()
		}
	})
}
