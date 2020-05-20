{{ define "RepeatedAssociations" }}
select
       B.id         as  blog_id,
       P.id         as  post_id,
       P.blog_id    as  post_blog_id,
       P.author_id  as  post_author_id,
       P.created_on as  post_created_on,
       P.section    as  post_section,
       P.subject    as  post_subject,
       P.draft      as  draft,
       P.body       as  post_body
from blog B
       left outer join post P on  B.id = P.blog_id
where B.id = 1
{{ end }}

{{ define "EmptyQuery" }}
{{ end }}

{{ define "InsertQueryAsExec" }}
select
       A.username          as  author_username,
       A.password          as  author_password,
       A.email             as  author_email,
       A.bio               as  author_bio,
       A.favourite_section as  author_favourite_section
from author A where id = 1
{{ end }}

{{ define "ExecAsQuery" }}
drop table if exists execasquery;
create table execasquery
(
  id                int
);
{{ end }}

{{ define "UnclaimedColumns" }}
select
       A.username          as  author_username,
       A.password          as  author_password
from author A where id = 1
{{ end }}

{{ define "MultipleRespForUnary" }}
select
       A.username          as  author_username,
       A.password          as  author_password,
       A.email             as  author_email,
       A.bio               as  author_bio,
       A.favourite_section as  author_favourite_section
from author A where id in (1,2)
{{ end }}

{{ define "NoRespForUnary" }}
select
       A.username          as  author_username,
       A.password          as  author_password,
       A.email             as  author_email,
       A.bio               as  author_bio,
       A.favourite_section as  author_favourite_section
from author A where id in (999)
{{ end }}

{{ define "RepeatedPrimative" }}
{{ end }}

{{ define "RepeatedEmpty" }}
select id as blog_id from blog limit 1
{{ end }}

{{ define "RepeatedTimestamp" }}
select
       B.id         as  blog_id,
       P.created_on as  post_created_on
from blog B
       left outer join post P on  B.id = P.blog_id
{{ end }}

{{ define "EmptyNestedField" }}
select id as  blog_id from blog limit 1
{{ end }}


{{ define "NoMatchingColumns" }}
select
       B.id as  blog_id,
       B.title as  blog_title
from blog B
{{ end }}

{{ define "AssociationInCollection" }}
select
        T.id as tag_id,
        P.id as post_id,
        B.id as blog_id
from tag T
        left outer join post_tag PT on PT.tag_id = T.id
        left outer join post P      on PT.post_id = P.id
        left outer join blog B      on P.blog_id = B.id
where T.id = 1
{{ end }}

{{ define "CollectionInAssociation" }}
select
        1 as dummy_var,
        B.id as  blog_id,
        B.title as  blog_title,
        A.id                as  author_id,
        A.username          as  author_username,
        A.password          as  author_password,
        A.email             as  author_email,
        A.bio               as  author_bio,
        A.favourite_section as  author_favourite_section
from blog B
       left outer join author A on  B.author_id = A.id
{{ end }}

{{ define "NullResoultsForSubmaps" }}
select
        PT.post_id as  post_id,
        P.id       as  id,
        C.id       as  comment_id,
        C.post_id  as  comment_post_id,
        C.comment  as  comment_text
from post_tag PT
       left outer join post P      on  PT.post_id = P.id
       left outer join comment C   on  P.id = C.post_id
where C.comment is null
{{ end }}

{{ define "SimpleEnum" }}
select
       A.username          as  author_username,
       A.password          as  author_password,
       A.email             as  author_email,
       A.bio               as  author_bio,
       A.favourite_section as  author_favourite_section
from author A where id = 3
{{ end }}

{{ define "NestedEnum" }}
select
        1 as id,
        2 as nested_id,
        'egg' as nested_enum
{{ end }}

{{ define "Blog" }}
select id from blog B order by id limit 1
{{ end }}

{{ define "Blogs" }}
select id from blog B order by id
{{ end }}

{{ define "BlogB" }}
{{ template "Blog" }}
{{ end }}

{{ define "BlogsB" }}
{{ template "Blogs" }}
{{ end }}

{{ define "BlogBF" }}
{{ template "Blog" }}
{{ end }}

{{ define "BlogsBF" }}
{{ template "Blogs" }}
{{ end }}

{{ define "BlogA" }}
{{ template "Blog" }}
{{ end }}

{{ define "BlogsA" }}
{{ template "Blogs" }}
{{ end }}

{{ define "BlogAF" }}
{{ template "Blog" }}
{{ end }}

{{ define "BlogsAF" }}
{{ template "Blogs" }}
{{ end }}

{{ define "BlogC" }}
{{ template "Blog" }}
{{ end }}

{{ define "BlogsC" }}
{{ template "Blogs" }}
{{ end }}

{{ define "BlogCF" }}
{{ template "Blog" }}
{{ end }}

{{ define "BlogsCF" }}
{{ template "Blogs" }}
{{ end }}

{{ define "CanceledUnaryContext" }}
select pg_sleep(15)
{{ end }}

{{ define "CanceledStreamContext" }}
select pg_sleep(15)
{{ end }}
