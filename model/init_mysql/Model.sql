# 用户信息表
# create table `User` (
#                         `id` bigint not null auto_increment unique comment 'PrimaryKey',
#                         `username` varchar(32) not null unique comment 'Name',
#                         `password` varchar(32) not null default 'tiktok' comment 'Password',
#                         `follow_count` bigint not null default 0 comment 'FollowCount',
#                         `follower_count` bigint not null default 0 comment 'FollowerCount',
#                         `is_follow` boolean not null default false comment 'IsFollow',
#                         `avatar` varchar(128) not null default '' comment 'Avatar',
#                         `background_image` varchar(128) not null default '' comment 'BackgroundImage',
#                         `signature` varchar(256) not null default 'gogogo' comment 'Signature',
#                         `total_favorited` bigint not null default 0 comment 'TotalFavorited',
#                         `work_count` bigint not null default 0 comment 'WorkCount',
#                         `favorite_count` bigint not null default 0 comment 'FavoriteCount',
#                         primary key (`id`)
# );
#
# # 视频信息表
# create table `Video` (
#     `id` bigint unique not null auto_increment comment 'Id',
#     `user_id` bigint unique not null comment 'UserID',
#     `play_url` varchar(128) not null default '' comment 'PlayUrl',
#     `cover_url` varchar(128) not null default '' comment 'CoverUrl',
#     `favorite_count` bigint not null default 0 comment 'FavoriteCount',
#     `comment_count` bigint not null default 0 comment 'CommentCount',
#     `is_favorite` boolean not null default false comment 'IsFavorite',
#     `title` varchar(256) not null comment 'Title',
#     `extension` varchar(32) not null default 'ext' comment 'Extension',
#     `created_time` timestamp not null default current_timestamp comment 'CreatedTime',
#     `updated_time` timestamp not null default current_timestamp on update current_timestamp comment 'UpdatedTime',
#     primary key (`id`),
#     foreign key (`user_id`) references User(`id`)
# )