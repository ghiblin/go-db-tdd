package docker_test

import (
	"fmt"

	godbtdd "github.com/ghiblin/go-db-tdd"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Repository", func() {
	var repo *godbtdd.Repository
	BeforeEach(func() {
		repo = &godbtdd.Repository{Db: Db}
		err := repo.Migrate()
		Expect(err).To(Succeed())

		_, err = repo.Create(&godbtdd.Blog{
			Title:   "post 1",
			Content: "hello",
			Tags:    []string{"tech", "fin"},
		})
		Expect(err).To(Succeed())
		_, err = repo.Create(&godbtdd.Blog{
			Title:   "post 2",
			Content: "world",
			Tags:    []string{"fin", "post"},
		})
		Expect(err).To(Succeed())
	})

	Context("Load", func() {
		It("Found", func() {
			blog, err := repo.Load(1)

			Expect(err).To(Succeed())
			Expect(blog.Title).To(Equal("post 1"))
			Expect(blog.Content).To(Equal("hello"))
			Expect(blog.Tags).To(Equal([]string{"tech", "fin"}))
		})

		It("Not Found", func() {
			_, err := repo.Load(999)

			Expect(err).To(HaveOccurred())
		})
	})

	It("ListAll", func() {
		blogs, err := repo.ListAll()
		Expect(err).To(Succeed())
		Expect(blogs).To(HaveLen(2))
	})

	It("List", func() {
		for i := 0; i < 20; i++ {
			_, err := repo.Create(&godbtdd.Blog{
				Title:   fmt.Sprintf("new post %v", i),
				Content: fmt.Sprintf("post %v content", i),
				Tags:    []string{"foo"},
			})
			Expect(err).To(Succeed())
		}
		l, err := repo.List(0, 10)
		Expect(err).To(Succeed())
		Expect(l).To(HaveLen(10))
	})

	Context("Save", func() {
		It("Create", func() {
			blog := &godbtdd.Blog{
				Title:   "post 3",
				Content: "hello",
				Tags:    []string{"foo"},
			}
			err := repo.Save(blog)
			Expect(err).To(Succeed())
			Expect(blog.ID).To(BeEquivalentTo(3))
		})

		It("Update", func() {
			blog, err := repo.Load(1)
			Expect(err).To(Succeed())

			blog.Title = "foo"
			err = repo.Save(blog)
			Expect(err).To(Succeed())

			blog, err = repo.Load(1)
			Expect(err).To(Succeed())
			Expect(blog.Title).To(Equal("foo"))
		})
	})

	It("Delete", func() {
		err := repo.Delete(1)
		Expect(err).To(Succeed())
		_, err = repo.Load(1)
		Expect(err).To(HaveOccurred())
	})

	DescribeTable("SearchByTitle",
		func(q string, found int) {
			l, err := repo.SearchByTitle(q, 0, 10)
			Expect(err).To(Succeed())
			Expect(l).To(HaveLen(found))
		},
		Entry("found", "post 1", 1),
		Entry("partial", "ost", 2),
		Entry("ignore case", "POST", 2),
		Entry("not found", "bar", 0),
	)

	DescribeTable("SearchByTag",
		func(q string, found int) {
			l, err := repo.SearchByTag(q, 0, 10)
			Expect(err).To(Succeed())
			Expect(l).To(HaveLen(found))
		},
		Entry("found", "tech", 1),
		Entry("not found", "foo", 0),
	)
})
