-- Copied from the corresponding version of dbx generated schema
CREATE TABLE accounting_raws (
	id bigserial NOT NULL,
	node_id bytea NOT NULL,
	interval_end_time timestamp with time zone NOT NULL,
	data_total double precision NOT NULL,
	data_type integer NOT NULL,
	created_at timestamp with time zone NOT NULL,
	PRIMARY KEY ( id )
);
CREATE TABLE accounting_rollups (
	id bigserial NOT NULL,
	node_id bytea NOT NULL,
	start_time timestamp with time zone NOT NULL,
	put_total bigint NOT NULL,
	get_total bigint NOT NULL,
	get_audit_total bigint NOT NULL,
	get_repair_total bigint NOT NULL,
	put_repair_total bigint NOT NULL,
	at_rest_total double precision NOT NULL,
	PRIMARY KEY ( id )
);
CREATE TABLE accounting_timestamps (
	name text NOT NULL,
	value timestamp with time zone NOT NULL,
	PRIMARY KEY ( name )
);
CREATE TABLE bwagreements (
	serialnum text NOT NULL,
	storage_node_id bytea NOT NULL,
	uplink_id bytea NOT NULL,
	action bigint NOT NULL,
	total bigint NOT NULL,
	created_at timestamp with time zone NOT NULL,
	expires_at timestamp with time zone NOT NULL,
	PRIMARY KEY ( serialnum )
);
CREATE TABLE certRecords (
	publickey bytea NOT NULL,
	id bytea NOT NULL,
	update_at timestamp with time zone NOT NULL,
	PRIMARY KEY ( id )
);
CREATE TABLE injuredsegments (
	id bigserial NOT NULL,
	info bytea NOT NULL,
	PRIMARY KEY ( id )
);
CREATE TABLE irreparabledbs (
	segmentpath bytea NOT NULL,
	segmentdetail bytea NOT NULL,
	pieces_lost_count bigint NOT NULL,
	seg_damaged_unix_sec bigint NOT NULL,
	repair_attempt_count bigint NOT NULL,
	PRIMARY KEY ( segmentpath )
);
CREATE TABLE nodes (
	id bytea NOT NULL,
	audit_success_count bigint NOT NULL,
	total_audit_count bigint NOT NULL,
	audit_success_ratio double precision NOT NULL,
	uptime_success_count bigint NOT NULL,
	total_uptime_count bigint NOT NULL,
	uptime_ratio double precision NOT NULL,
	created_at timestamp with time zone NOT NULL,
	updated_at timestamp with time zone NOT NULL,
	wallet text NOT NULL,
	email text NOT NULL,
	PRIMARY KEY ( id )
);
CREATE TABLE overlay_cache_nodes (
	node_id bytea NOT NULL,
	node_type integer NOT NULL,
	address text NOT NULL,
	protocol integer NOT NULL,
	operator_email text NOT NULL,
	operator_wallet text NOT NULL,
	free_bandwidth bigint NOT NULL,
	free_disk bigint NOT NULL,
	latency_90 bigint NOT NULL,
	audit_success_ratio double precision NOT NULL,
	audit_uptime_ratio double precision NOT NULL,
	audit_count bigint NOT NULL,
	audit_success_count bigint NOT NULL,
	uptime_count bigint NOT NULL,
	uptime_success_count bigint NOT NULL,
	PRIMARY KEY ( node_id ),
	UNIQUE ( node_id )
);
CREATE TABLE projects (
	id bytea NOT NULL,
	name text NOT NULL,
	description text NOT NULL,
	created_at timestamp with time zone NOT NULL,
	PRIMARY KEY ( id )
);
CREATE TABLE users (
	id bytea NOT NULL,
	first_name text NOT NULL,
	last_name text NOT NULL,
	email text NOT NULL,
	password_hash bytea NOT NULL,
	status integer NOT NULL,
	created_at timestamp with time zone NOT NULL,
	PRIMARY KEY ( id )
);
CREATE TABLE api_keys (
	id bytea NOT NULL,
	project_id bytea NOT NULL REFERENCES projects( id ) ON DELETE CASCADE,
	key bytea NOT NULL,
	name text NOT NULL,
	created_at timestamp with time zone NOT NULL,
	PRIMARY KEY ( id ),
	UNIQUE ( key ),
	UNIQUE ( name, project_id )
);
CREATE TABLE project_members (
	member_id bytea NOT NULL REFERENCES users( id ) ON DELETE CASCADE,
	project_id bytea NOT NULL REFERENCES projects( id ) ON DELETE CASCADE,
	created_at timestamp with time zone NOT NULL,
	PRIMARY KEY ( member_id, project_id )
);

INSERT INTO "accounting_raws" VALUES (1, E'\\3510\\323\\225"~\\036<\\342\\330m\\0253Jhr\\246\\233K\\246#\\2303\\351\\256\\275j\\212UM\\362\\207', '2019-02-14 08:16:57.812849+00', 1000, 0, '2019-02-14 08:16:57.844849+00');

INSERT INTO "accounting_rollups"("id", "node_id", "start_time", "put_total", "get_total", "get_audit_total", "get_repair_total", "put_repair_total", "at_rest_total") VALUES (1, E'\\367M\\177\\251]t/\\022\\256\\214\\265\\025\\224\\204:\\217\\212\\0102<\\321\\374\\020&\\271Qc\\325\\261\\354\\246\\233'::bytea, '2019-02-09 00:00:00+00', 1000, 2000, 3000, 4000, 0, 5000);

INSERT INTO "accounting_timestamps" VALUES ('LastAtRestTally', '0001-01-01 00:00:00+00');
INSERT INTO "accounting_timestamps" VALUES ('LastRollup', '0001-01-01 00:00:00+00');
INSERT INTO "accounting_timestamps" VALUES ('LastBandwidthTally', '0001-01-01 00:00:00+00');

INSERT INTO "nodes" VALUES (E'\\006\\223\\250R\\221\\005\\365\\377v>0\\266\\365\\216\\255?\\347\\244\\371?2\\264\\262\\230\\007<\\001\\262\\263\\237\\247n', 0, 0, 0, 3, 3, 1, '2019-02-14 08:07:31.028103+00', '2019-02-14 08:07:31.108963+00', '', '');
INSERT INTO "overlay_cache_nodes" VALUES (E'\\006\\223\\250R\\221\\005\\365\\377v>0\\266\\365\\216\\255?\\347\\244\\371?2\\264\\262\\230\\007<\\001\\262\\263\\237\\247n', 4, '127.0.0.1:55518', 0, 'bootstrap@mail.test', '0x0000000000000000000000000000000000000000', -1, -1, 0, 0, 1, 0, 0, 2, 2);

INSERT INTO "projects"("id", "name", "description", "created_at") VALUES (E'\\022\\217/\\014\\376!K\\023\\276\\031\\311}m\\236\\205\\300'::bytea, 'ProjectName', 'projects description', '2019-02-14 08:28:24.254934+00');
INSERT INTO "api_keys"("id", "project_id", "key", "name", "created_at") VALUES (E'\\334/\\302;\\225\\355O\\323\\276f\\247\\354/6\\241\\033'::bytea, E'\\022\\217/\\014\\376!K\\023\\276\\031\\311}m\\236\\205\\300'::bytea, E'\\000]\\326N \\343\\270L\\327\\027\\337\\242\\240\\322mOl\\0318\\251.P I'::bytea, 'key 2', '2019-02-14 08:28:24.267934+00');

INSERT INTO "users"("id", "first_name", "last_name", "email", "password_hash", "status", "created_at") VALUES (E'\\363\\311\\033w\\222\\303Ci\\265\\343U\\303\\312\\204",'::bytea, 'Noahson', 'William', '1email1@mail.test', E'some_readable_hash'::bytea, 1, '2019-02-14 08:28:24.614594+00');
INSERT INTO "projects"("id", "name", "description", "created_at") VALUES (E'\\363\\342\\363\\371>+F\\256\\263\\300\\273|\\342N\\347\\014'::bytea, 'projName1', 'Test project 1', '2019-02-14 08:28:24.636949+00');
INSERT INTO "project_members"("member_id", "project_id", "created_at") VALUES (E'\\363\\311\\033w\\222\\303Ci\\265\\343U\\303\\312\\204",'::bytea, E'\\363\\342\\363\\371>+F\\256\\263\\300\\273|\\342N\\347\\014'::bytea, '2019-02-14 08:28:24.677953+00');

INSERT INTO "bwagreements"("serialnum", "storage_node_id", "action", "total", "created_at", "expires_at", "uplink_id") VALUES ('8fc0ceaa-984c-4d52-bcf4-b5429e1e35e812FpiifDbcJkePa12jxjDEutKrfLmwzT7sz2jfVwpYqgtM8B74c', E'\\245Z[/\\333\\022\\011\\001\\036\\003\\204\\005\\032.\\206\\333E\\261\\342\\227=y,}aRaH6\\240\\370\\000'::bytea, 1, 666, '2019-02-14 15:09:54.420181+00', '2019-02-14 16:09:54+00', E'\\253Z+\\374eFm\\245$\\036\\206\\335\\247\\263\\350x\\\\\\304+\\364\\343\\364+\\276fIJQ\\361\\014\\232\\000'::bytea);
INSERT INTO "irreparabledbs" ("segmentpath", "segmentdetail", "pieces_lost_count", "seg_damaged_unix_sec", "repair_attempt_count") VALUES ('\x49616d5365676d656e746b6579696e666f30', '\x49616d5365676d656e7464657461696c696e666f30', 10, 1550159554, 10);
INSERT INTO "injuredsegments" ("id", "info") VALUES (1, '\x0a0130120100');

INSERT INTO "certrecords" VALUES (E'0Y0\\023\\006\\007*\\206H\\316=\\002\\001\\006\\010*\\206H\\316=\\003\\001\\007\\003B\\000\\004\\360\\267\\227\\377\\253u\\222\\337Y\\324C:GQ\\010\\277v\\010\\315D\\271\\333\\337.\\203\\023=C\\343\\014T%6\\027\\362?\\214\\326\\017U\\334\\000\\260\\224\\260J\\221\\304\\331F\\304\\221\\236zF,\\325\\326l\\215\\306\\365\\200\\022', E'L\\301|\\200\\247}F|1\\320\\232\\037n\\335\\241\\206\\244\\242\\207\\204.\\253\\357\\326\\352\\033Dt\\202`\\022\\325', '2019-02-14 08:07:31.335028+00');
