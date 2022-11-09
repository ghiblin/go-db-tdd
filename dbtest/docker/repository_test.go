package docker_test

import (
	godbtdd "github.com/ghiblin/go-db-tdd"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Repository", func() {
	var repo *godbtdd.Repository
	BeforeEach(func() {
		repo = &godbtdd.Repository{Db: Db}
		err := repo.Migrate()
		Ω(err).To(Succeed())

		sampleData := &godbtdd.Blog{
			Title:   "post",
			Content: "hello",
			Tags:    []string{"a", "b"},
		}
		_, err = repo.Create(sampleData)
		Ω(err).To(Succeed())
	})

	Context("Load", func() {
		It("Found", func() {
			blog, err := repo.Load(1)

			Ω(err).To(Succeed())
			Ω(blog.Title).To(Equal("post"))
			Ω(blog.Content).To(Equal("hello"))
			Ω(blog.Tags).To(Equal([]string{"a", "b"}))
		})
	})
})
