package echopen

import (
	v310 "github.com/richjyoung/echopen/openapi/v3.1.0"
)

type SpecFilterFunc = func(s *v310.Specification) *v310.Specification

func IncludeTags(tags ...string) SpecFilterFunc {
	return func(s *v310.Specification) *v310.Specification {
		newTags := []*v310.Tag{}
		for _, tag := range s.Tags {
			for _, included := range tags {
				if tag.Name == included {
					newTags = append(newTags, tag)
				}
			}
		}

		s.Tags = newTags

		for name, path := range s.Paths {
			if path.Value != nil {
				if path.Value.Delete != nil && len(path.Value.Delete.Tags) > 0 {
					path.Value.Delete.Tags = filterStringSliceIncludes(tags, path.Value.Delete.Tags)
					if len(path.Value.Delete.Tags) == 0 {
						path.Value.Delete = nil
					}
				}

				if path.Value.Get != nil && len(path.Value.Get.Tags) > 0 {
					path.Value.Get.Tags = filterStringSliceIncludes(tags, path.Value.Get.Tags)
					if len(path.Value.Get.Tags) == 0 {
						path.Value.Get = nil
					}
				}

				if path.Value.Head != nil && len(path.Value.Head.Tags) > 0 {
					path.Value.Head.Tags = filterStringSliceIncludes(tags, path.Value.Head.Tags)
					if len(path.Value.Head.Tags) == 0 {
						path.Value.Head = nil
					}
				}

				if path.Value.Options != nil && len(path.Value.Options.Tags) > 0 {
					path.Value.Options.Tags = filterStringSliceIncludes(tags, path.Value.Options.Tags)
					if len(path.Value.Options.Tags) == 0 {
						path.Value.Options = nil
					}
				}

				if path.Value.Patch != nil && len(path.Value.Patch.Tags) > 0 {
					path.Value.Patch.Tags = filterStringSliceIncludes(tags, path.Value.Patch.Tags)
					if len(path.Value.Patch.Tags) == 0 {
						path.Value.Patch = nil
					}
				}

				if path.Value.Post != nil && len(path.Value.Post.Tags) > 0 {
					path.Value.Post.Tags = filterStringSliceIncludes(tags, path.Value.Post.Tags)
					if len(path.Value.Post.Tags) == 0 {
						path.Value.Post = nil
					}
				}

				if path.Value.Put != nil && len(path.Value.Put.Tags) > 0 {
					path.Value.Put.Tags = filterStringSliceIncludes(tags, path.Value.Put.Tags)
					if len(path.Value.Put.Tags) == 0 {
						path.Value.Put = nil
					}
				}

				if path.Value.Trace != nil && len(path.Value.Trace.Tags) > 0 {
					path.Value.Trace.Tags = filterStringSliceIncludes(tags, path.Value.Trace.Tags)
					if len(path.Value.Trace.Tags) == 0 {
						path.Value.Trace = nil
					}
				}
			}

			if path.Value.Delete == nil &&
				path.Value.Get == nil &&
				path.Value.Head == nil &&
				path.Value.Options == nil &&
				path.Value.Patch == nil &&
				path.Value.Post == nil &&
				path.Value.Put == nil &&
				path.Value.Trace == nil {

				delete(s.Paths, name)
			}
		}

		return s
	}
}

func ExcludeTags(tags ...string) SpecFilterFunc {
	return func(s *v310.Specification) *v310.Specification {
		newTags := []*v310.Tag{}
		for _, tag := range s.Tags {
			for _, excluded := range tags {
				if tag.Name == excluded {
					continue
				}
				newTags = append(newTags, tag)
			}
		}

		s.Tags = newTags

		for name, path := range s.Paths {
			if path.Value != nil {

				if path.Value.Delete != nil && len(path.Value.Delete.Tags) > 0 {
					path.Value.Delete.Tags = filterStringSliceExcludes(tags, path.Value.Delete.Tags)
					if len(path.Value.Delete.Tags) == 0 {
						path.Value.Delete = nil
					}
				}

				if path.Value.Get != nil && len(path.Value.Get.Tags) > 0 {
					path.Value.Get.Tags = filterStringSliceExcludes(tags, path.Value.Get.Tags)
					if len(path.Value.Get.Tags) == 0 {
						path.Value.Get = nil
					}
				}

				if path.Value.Head != nil && len(path.Value.Head.Tags) > 0 {
					path.Value.Head.Tags = filterStringSliceExcludes(tags, path.Value.Head.Tags)
					if len(path.Value.Head.Tags) == 0 {
						path.Value.Head = nil
					}
				}

				if path.Value.Options != nil && len(path.Value.Options.Tags) > 0 {
					path.Value.Options.Tags = filterStringSliceExcludes(tags, path.Value.Options.Tags)
					if len(path.Value.Options.Tags) == 0 {
						path.Value.Options = nil
					}
				}

				if path.Value.Patch != nil && len(path.Value.Patch.Tags) > 0 {
					path.Value.Patch.Tags = filterStringSliceExcludes(tags, path.Value.Patch.Tags)
					if len(path.Value.Patch.Tags) == 0 {
						path.Value.Patch = nil
					}
				}

				if path.Value.Post != nil && len(path.Value.Post.Tags) > 0 {
					path.Value.Post.Tags = filterStringSliceExcludes(tags, path.Value.Post.Tags)
					if len(path.Value.Post.Tags) == 0 {
						path.Value.Post = nil
					}
				}

				if path.Value.Put != nil && len(path.Value.Put.Tags) > 0 {
					path.Value.Put.Tags = filterStringSliceExcludes(tags, path.Value.Put.Tags)
					if len(path.Value.Put.Tags) == 0 {
						path.Value.Put = nil
					}
				}

				if path.Value.Trace != nil && len(path.Value.Trace.Tags) > 0 {
					path.Value.Trace.Tags = filterStringSliceExcludes(tags, path.Value.Trace.Tags)
					if len(path.Value.Trace.Tags) == 0 {
						path.Value.Trace = nil
					}
				}
			}

			if path.Value.Delete == nil &&
				path.Value.Get == nil &&
				path.Value.Head == nil &&
				path.Value.Options == nil &&
				path.Value.Patch == nil &&
				path.Value.Post == nil &&
				path.Value.Put == nil &&
				path.Value.Trace == nil {

				delete(s.Paths, name)
			}
		}

		return s
	}
}

func filterStringSliceIncludes(include []string, slice []string) []string {
	s := []string{}

	for _, i := range include {
		for _, v := range slice {
			if i == v {
				s = append(s, v)
			}
		}
	}

	return s
}

func filterStringSliceExcludes(excludes []string, slice []string) []string {
	s := []string{}

	for _, v := range slice {
		for _, i := range excludes {
			if i == v {
				continue
			}
			s = append(s, v)
		}
	}

	return s
}
