create database sakilogy;

--- enjoy hard coding password
create user 'sakilogy'@'localhost' identified by '@k052a9';
grant all privileges on sakilogy.* to 'sakilogy'@'localhost';
flush privileges;

use sakilogy;

create table users (
    user_id int unsigned primary key not null auto_increment,
    username varchar(16) not null unique,
    password char(44) not null,
    level int(2) unsigned not null default 0,
    pt int(4) unsigned not null default 0,
    rating float(8,3) not null default 1500.0
);

create table girls (
    girl_id int primary key,
    level int(2) unsigned not null default 0,
    pt int(4) unsigned not null default 0,
    rating float(8,3) not null default 1500.0
);

-- a sample girl list, written at v0.7.6. add more if necessary.
insert into girls (girl_id) values
	(0),
	(710113), (710114), (710115),
	(712411), (712412), (712413),
	(712611), (712613),
	(712714), (712715),
	(712915),
	(713311), (713314),
	(713301),
	(713811), (713815),
	(714915),
	(715712),
	(990001), (990002);

create table user_girl (
    user_id int unsigned not null,
    girl_id int not null,
    rank1 int unsigned not null default 0,
    rank2 int unsigned not null default 0,
    rank3 int unsigned not null default 0,
    rank4 int unsigned not null default 0,
    primary key (user_id,girl_id),
    foreign key (user_id) references users(user_id),
    foreign key (girl_id) references girls(girl_id)
);

alter table user_girl
    add column play int unsigned
    generated always as (rank1+rank2+rank3+rank4) virtual;

alter table user_girl add column (
    avg_point float(8,1) not null default 0,
	a_top int unsigned not null default 0,
	a_last int unsigned not null default 0,
    round int unsigned not null default 0,
    win int unsigned not null default 0,
    gun int unsigned not null default 0,
    bark int unsigned not null default 0,
    riichi int unsigned not null default 0,
    win_point float(8,1) not null default 0,
    gun_point float(8,1) not null default 0,
    bark_point float(8,1) not null default 0,
    riichi_point float(8,1) not null default 0,
    ready int unsigned not null default 0,
    ready_turn float(4, 2) not null default 0,
    win_turn float(4, 2) not null default 0
);

alter table user_girl add column (
	yaku_rci int unsigned not null default 0,
	yaku_ipt int unsigned not null default 0,
	yaku_tmo int unsigned not null default 0,
	yaku_tny int unsigned not null default 0,
	yaku_pnf int unsigned not null default 0,
	yaku_y1y int unsigned not null default 0,
	yaku_y2y int unsigned not null default 0,
	yaku_y3y int unsigned not null default 0,
	yaku_jk1 int unsigned not null default 0,
	yaku_jk2 int unsigned not null default 0,
	yaku_jk3 int unsigned not null default 0,
	yaku_jk4 int unsigned not null default 0,
	yaku_bk1 int unsigned not null default 0,
	yaku_bk2 int unsigned not null default 0,
	yaku_bk3 int unsigned not null default 0,
	yaku_bk4 int unsigned not null default 0,
	yaku_ipk int unsigned not null default 0,
	yaku_rns int unsigned not null default 0,
	yaku_hai int unsigned not null default 0,
	yaku_hou int unsigned not null default 0,
	yaku_ckn int unsigned not null default 0,
	yaku_ss1 int unsigned not null default 0,
	yaku_it1 int unsigned not null default 0,
	yaku_ct1 int unsigned not null default 0,
	yaku_wri int unsigned not null default 0,
	yaku_ss2 int unsigned not null default 0,
	yaku_it2 int unsigned not null default 0,
	yaku_ct2 int unsigned not null default 0,
	yaku_toi int unsigned not null default 0,
	yaku_ctt int unsigned not null default 0,
	yaku_sak int unsigned not null default 0,
	yaku_skt int unsigned not null default 0,
	yaku_stk int unsigned not null default 0,
	yaku_hrt int unsigned not null default 0,
	yaku_s3g int unsigned not null default 0,
	yaku_h1t int unsigned not null default 0,
	yaku_jc2 int unsigned not null default 0,
	yaku_mnh int unsigned not null default 0,
	yaku_jc3 int unsigned not null default 0,
	yaku_rpk int unsigned not null default 0,
	yaku_c1t int unsigned not null default 0,
	yaku_mnc int unsigned not null default 0,
	yaku_x13 int unsigned not null default 0,
	yaku_xd3 int unsigned not null default 0,
	yaku_x4a int unsigned not null default 0,
	yaku_xt1 int unsigned not null default 0,
	yaku_xs4 int unsigned not null default 0,
	yaku_xd4 int unsigned not null default 0,
	yaku_xcr int unsigned not null default 0,
	yaku_xr1 int unsigned not null default 0,
	yaku_xth int unsigned not null default 0,
	yaku_xch int unsigned not null default 0,
	yaku_x4k int unsigned not null default 0,
	yaku_x9r int unsigned not null default 0,
	yaku_w13 int unsigned not null default 0,
	yaku_w4a int unsigned not null default 0,
	yaku_w9r int unsigned not null default 0,
	kzeykm int unsigned not null default 0,
	han_rci float(4,2) not null default 0,
	han_ipt float(4,2) not null default 0,
	han_tmo float(4,2) not null default 0,
	han_tny float(4,2) not null default 0,
	han_pnf float(4,2) not null default 0,
	han_y1y float(4,2) not null default 0,
	han_y2y float(4,2) not null default 0,
	han_y3y float(4,2) not null default 0,
	han_jk1 float(4,2) not null default 0,
	han_jk2 float(4,2) not null default 0,
	han_jk3 float(4,2) not null default 0,
	han_jk4 float(4,2) not null default 0,
	han_bk1 float(4,2) not null default 0,
	han_bk2 float(4,2) not null default 0,
	han_bk3 float(4,2) not null default 0,
	han_bk4 float(4,2) not null default 0,
	han_ipk float(4,2) not null default 0,
	han_rns float(4,2) not null default 0,
	han_hai float(4,2) not null default 0,
	han_hou float(4,2) not null default 0,
	han_ckn float(4,2) not null default 0,
	han_ss1 float(4,2) not null default 0,
	han_it1 float(4,2) not null default 0,
	han_ct1 float(4,2) not null default 0,
	han_wri float(4,2) not null default 0,
	han_ss2 float(4,2) not null default 0,
	han_it2 float(4,2) not null default 0,
	han_ct2 float(4,2) not null default 0,
	han_toi float(4,2) not null default 0,
	han_ctt float(4,2) not null default 0,
	han_sak float(4,2) not null default 0,
	han_skt float(4,2) not null default 0,
	han_stk float(4,2) not null default 0,
	han_hrt float(4,2) not null default 0,
	han_s3g float(4,2) not null default 0,
	han_h1t float(4,2) not null default 0,
	han_jc2 float(4,2) not null default 0,
	han_mnh float(4,2) not null default 0,
	han_jc3 float(4,2) not null default 0,
	han_rpk float(4,2) not null default 0,
	han_c1t float(4,2) not null default 0,
	han_mnc float(4,2) not null default 0,
	yaku_dora int unsigned not null default 0,
	yaku_uradora int unsigned not null default 0,
	yaku_akadora int unsigned not null default 0,
	yaku_kandora int unsigned not null default 0,
	yaku_kanuradora int unsigned not null default 0
);

create table replays (
    replay_id int unsigned primary key not null auto_increment,
	content mediumtext not null
);

create table user_replay (
    user_id int unsigned not null,
    replay_id int unsigned not null,
    primary key (user_id, replay_id),
    foreign key (user_id) references users(user_id),
    foreign key (replay_id) references replays(replay_id)
);
