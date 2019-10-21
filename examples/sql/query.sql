{{ define "SelectBlog" }}
select
        id,
        title,
        author_id
from blog
where id = {{ .Id }} limit 1
{{ end }}


{{ define "SelectBlogs" }}
select
        id,
        title,
        author_id
from blog
        where id in (
        {{ .Ids | join " , " }}
        ) or
        title in (
        {{ .Titles | squoteall | join " , " }}
        )
{{ end }}


{{ define "SelectDetailedBlog" }}
        {{template "DetailedBlog" }}
        {{ if .Id }}
                where B.id = {{ .Id }}
        {{ else if .AuthorId }}
                where B.author_id = {{ .AuthorId }}
        {{ end }}
{{ end }}


{{ define "SelectDetailedBlogs" }}
        {{ template "DetailedBlog" }}
        where B.id in (
        {{ .Ids | join " , " }}
        ) or
        B.title in (
        {{ .Titles | squoteall | join " , " }}
        )
{{ end }}

{{ define "DetailedBlog" }}
select
        B.id                as  blog_id,
        B.title             as  blog_title,
        A.id                as  author_id,
        A.username          as  author_username,
        A.password          as  author_password,
        A.email             as  author_email,
        A.bio               as  author_bio,
        A.favourite_section as  author_favourite_section,
        P.id                as  post_id,
        P.blog_id           as  post_blog_id,
        P.author_id         as  post_author_id,
        P.created_on        as  post_created_on,
        P.section           as  post_section,
        P.subject           as  post_subject,
        P.draft             as  draft,
        P.body              as  post_body,
        C.id                as  comment_id,
        C.post_id           as  comment_post_id,
        C.comment           as  comment_text,
        T.id                as  tag_id,
        T.name              as  tag_name
from blog B
        left outer join author A    on  B.author_id = A.id
        left outer join post P      on  B.id = P.blog_id
        left outer join comment C   on  P.id = C.post_id
        left outer join post_tag PT on  PT.post_id = P.id
        left outer join tag T       on  PT.tag_id = T.id
{{ end }}


{{define "InsertAuthor" }}
INSERT INTO author
VALUES (
 {{ .Id }},
 {{ .Username | quote }},
 {{ .Password | quote }},
 {{ .Email | quote }},
 {{ .Bio | quote }},
 {{ .FavouriteSection | quote }}
);
{{end}}
