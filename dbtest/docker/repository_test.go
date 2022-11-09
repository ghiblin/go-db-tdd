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

		_, err = repo.Create(&godbtdd.Blog{
			Title:   "post 1",
			Content: "hello",
			Tags:    []string{"a", "b"},
		})
		Ω(err).To(Succeed())
		_, err = repo.Create(&godbtdd.Blog{
			Title:   "post 2",
			Content: "world",
			Tags:    []string{"b", "c"},
		})
		Ω(err).To(Succeed())
	})

	Context("Load", func() {
		It("Found", func() {
			blog, err := repo.Load(1)

			Ω(err).To(Succeed())
			Ω(blog.Title).To(Equal("post 1"))
			Ω(blog.Content).To(Equal("hello"))
			Ω(blog.Tags).To(Equal([]string{"a", "b"}))
		})

		It("Not Found", func() {
			_, err := repo.Load(999)

			Ω(err).To(HaveOccurred())
		})
	})

	It("ListAll", func() {
		blogs, err := repo.ListAll()
		Ω(err).To(Succeed())
		Ω(blogs).To(HaveLen(2))
	})
})
