package sql_scripts

func TableSchema() string {
	return `
create table if not exists public.Users
(
    UserID      serial
        constraint Users_pk
            primary key,
    Username    varchar(26)  not null,
    Email       varchar(31)  not null,
    IsAdmin     boolean,
    DisplayName varchar(26),
    UserSub     varchar(40) not null
);

create table if not exists public.Friends
(
    FriendsID serial
        constraint Friends_pk
            primary key,
    User1ID   serial not null
        constraint Friends_users_userid_fk_1
            references public.users,
    User2ID   serial not null
        constraint Friends_users_userid_fk_2
            references public.users
);

create table if not exists public.FriendRequests
(
    FriendRequestID serial
        constraint FriendRequests_pk
            primary key,
    SenderUserID    serial not null
        constraint FriendRequests_users_userid_fk
            references public.users,
    ReceiverUserID  serial not null
        constraint FriendRequests_users_userid_fk_2
            references public.users
);

create table if not exists public.UserStats
(
    UserStatsID          serial
        constraint UserStats_pk
            primary key,
    GamesPlayed          integer default 0 not null,
    GamesWon             integer default 0 not null,
    RoundsPlayed         integer default 0 not null,
    RoundsWon            integer default 0 not null,
    TotalPointsAsBidder  integer default 0 not null,
    TotalRoundsAsBidder  integer default 0 not null,
    TotalPointsAsSupport integer default 0 not null,
    TotalRoundsAsSupport integer default 0 not null,
    TotalPointsAsCounter integer default 0 not null,
    TotalRoundsAsCounter integer default 0 not null,
    TimesWinningBidTotal integer default 0 not null,
    TimesCallingSuit     integer default 0 not null,
    TimesCallingNil      integer default 0 not null,
    TimesCallingSplash   integer default 0 not null,
    TimesCallingPlunge   integer default 0 not null,
    TimesCallingSevens   integer default 0 not null,
    TimesCallingDelve    integer default 0 not null,
    userid               integer           not null
        constraint userstats_users_userid_fk
            references public.users
);

create table if not exists public.MatchArchive
(
    MatchID      serial
        constraint MatchArchive_pk
            primary key,
    MatchName    varchar(41)       not null,
    MatchPrivacy smallint          not null,
    Rules        bit varying       not null,
    TotalRounds  integer default 0 not null,
    WinningTeam  integer default 1 not null,
    Team1Marks   integer default 0 not null,
    Team2Marks   integer default 0 not null,
    Team1Player1 serial            not null
        constraint MatchArchive_users_userid_fk
            references public.users,
    Team1Player2 serial            not null
        constraint MatchArchive_users_userid_fk_2
            references public.users,
    Team2Player1 serial            not null
        constraint MatchArchive_users_userid_fk_3
            references public.users,
    Team2Player2 serial            not null
        constraint MatchArchive_users_userid_fk_4
            references public.users
);

create table if not exists public.RoundArchive
(
    RoundID       serial
        constraint RoundArchive_pk
            primary key,
    MatchID       serial      not null
        constraint RoundArchive_matcharchive_matchid_fk
            references public.matcharchive,
    RoundRules    bit varying not null,
    Team1Score    integer     not null,
    Team2Score    integer     not null,
    RoundActivity varchar     not null
);

create table if not exists public.ChatMessageArchive
(
    ChatMessageID    serial
        constraint ChatMessageArchive_pk
            primary key,
    MatchID          serial       not null
        constraint ChatMessageArchive_matcharchive_matchid_fk
            references public.matcharchive,
    UserID           serial       not null
        constraint ChatMessageArchive_users_userid_fk
            references public.users,
    MessageTimestamp timestamp    not null,
    Message          varchar(256) not null
);

create table if not exists Logs
(
    LogID       serial
        constraint logs_pk
            primary key,
    Created     timestamp default now() not null,
    ServiceName varchar(30)             not null,
    Message     varchar                 not null
);
`
}
