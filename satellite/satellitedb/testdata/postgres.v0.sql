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
	data bytea NOT NULL,
	storage_node bytea NOT NULL,
	action bigint NOT NULL,
	total bigint NOT NULL,
	created_at timestamp with time zone NOT NULL,
	expires_at timestamp with time zone NOT NULL,
	PRIMARY KEY ( serialnum )
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
	email text,
	password_hash bytea NOT NULL,
	created_at timestamp with time zone NOT NULL,
	PRIMARY KEY ( id ),
	UNIQUE ( email )
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
CREATE TABLE bucket_infos (
	project_id bytea NOT NULL REFERENCES projects( id ) ON DELETE CASCADE,
	name text NOT NULL,
	created_at timestamp with time zone NOT NULL,
	PRIMARY KEY ( name )
);
CREATE TABLE project_members (
	member_id bytea NOT NULL REFERENCES users( id ) ON DELETE CASCADE,
	project_id bytea NOT NULL REFERENCES projects( id ) ON DELETE CASCADE,
	created_at timestamp with time zone NOT NULL,
	PRIMARY KEY ( member_id, project_id )
); 

-- NEW DATA --

INSERT INTO "accounting_raws" VALUES (1, E'\\3510\\323\\225"~\\036<\\342\\330m\\0253Jhr\\246\\233K\\246#\\2303\\351\\256\\275j\\212UM\\362\\207', '2019-02-14 08:16:57.812849+00', 1000, 0, '2019-02-14 08:16:57.844849+00');

INSERT INTO "accounting_rollups"("id", "node_id", "start_time", "put_total", "get_total", "get_audit_total", "get_repair_total", "put_repair_total", "at_rest_total") VALUES (1, E'\\367M\\177\\251]t/\\022\\256\\214\\265\\025\\224\\204:\\217\\212\\0102<\\321\\374\\020&\\271Qc\\325\\261\\354\\246\\233'::bytea, '2019-02-09 00:00:00+00', 1000, 2000, 3000, 4000, 0, 5000);

INSERT INTO "accounting_timestamps" VALUES ('LastAtRestTally', '0001-01-01 00:00:00+00');
INSERT INTO "accounting_timestamps" VALUES ('LastRollup', '0001-01-01 00:00:00+00');
INSERT INTO "accounting_timestamps" VALUES ('LastBandwidthTally', '0001-01-01 00:00:00+00');

INSERT INTO "nodes" VALUES (E'\\006\\223\\250R\\221\\005\\365\\377v>0\\266\\365\\216\\255?\\347\\244\\371?2\\264\\262\\230\\007<\\001\\262\\263\\237\\247n', 0, 0, 0, 3, 3, 1, '2019-02-14 08:07:31.028103+00', '2019-02-14 08:07:31.108963+00');
INSERT INTO "overlay_cache_nodes" VALUES (E'\\006\\223\\250R\\221\\005\\365\\377v>0\\266\\365\\216\\255?\\347\\244\\371?2\\264\\262\\230\\007<\\001\\262\\263\\237\\247n', 4, '127.0.0.1:55518', 0, 'bootstrap@mail.test', '0x0000000000000000000000000000000000000000', -1, -1, 0, 0, 1, 0, 0, 2, 2);

INSERT INTO "projects"("id", "name", "description", "created_at") VALUES (E'\\022\\217/\\014\\376!K\\023\\276\\031\\311}m\\236\\205\\300'::bytea, 'ProjectName', 'projects description', '2019-02-14 08:28:24.254934+00');
INSERT INTO "api_keys"("id", "project_id", "key", "name", "created_at") VALUES (E'\\334/\\302;\\225\\355O\\323\\276f\\247\\354/6\\241\\033'::bytea, E'\\022\\217/\\014\\376!K\\023\\276\\031\\311}m\\236\\205\\300'::bytea, E'\\000]\\326N \\343\\270L\\327\\027\\337\\242\\240\\322mOl\\0318\\251.P I'::bytea, 'key 2', '2019-02-14 08:28:24.267934+00');

INSERT INTO "users"("id", "first_name", "last_name", "email", "password_hash", "created_at") VALUES (E'\\363\\311\\033w\\222\\303Ci\\265\\343U\\303\\312\\204",'::bytea, 'Noahson', 'William', '1email1@mail.test', E'some_readable_hash'::bytea, '2019-02-14 08:28:24.614594+00');
INSERT INTO "projects"("id", "name", "description", "created_at") VALUES (E'\\363\\342\\363\\371>+F\\256\\263\\300\\273|\\342N\\347\\014'::bytea, 'projName1', 'Test project 1', '2019-02-14 08:28:24.636949+00');
INSERT INTO "project_members"("member_id", "project_id", "created_at") VALUES (E'\\363\\311\\033w\\222\\303Ci\\265\\343U\\303\\312\\204",'::bytea, E'\\363\\342\\363\\371>+F\\256\\263\\300\\273|\\342N\\347\\014'::bytea, '2019-02-14 08:28:24.677953+00');

INSERT INTO "bwagreements" ("serialnum", "data", "storage_node", "action", "total", "created_at", "expires_at") VALUES ('8fc0ceaa-984c-4d52-bcf4-b5429e1e35e812FpiifDbcJkePa12jxjDEutKrfLmwzT7sz2jfVwpYqgtM8B74c', '\x0a84070a206a4f331aeb1089097730bb01a3ac80851b9b04daa29fa41ae65796620253ec001220ab5a2bfc65466da5241e86dda7b3e8785cc42bf4e3f42bbe66494a51f10c9a0020d2a596e3052a2438666330636561612d393834632d346435322d626366342d623534323965316533356538300138c28996e30542e5023082016130820107a0030201020210656e4090e70c3c011edfc54352d53add300a06082a8648ce3d0403023010310e300c060355040a130553746f726a3022180f30303031303130313030303030305a180f30303031303130313030303030305a3010310e300c060355040a130553746f726a3059301306072a8648ce3d020106082a8648ce3d03010703420004189f0e56bd7bd36fa172ae7aed8aee9a35064397dad743bad61bcedcdd6534b3ccb7dde8b8331dd51266bd28cce026b13594c8ca225ef0dde82045aefbb72718a33f303d300e0603551d0f0101ff0404030205a0301d0603551d250416301406082b0601050507030106082b06010505070302300c0603551d130101ff04023000300a06082a8648ce3d0403020348003045022100c5ee1c6c2eb24a32df27b3659705266017f27e65fe07134095f73ea8b1f064b70220305b166d1e718f0894c6b56fca1c09cc9912a7e91fbe5371a81006ee66533eb342df023082015b30820101a003020102021100b50f5053d4d52fa0935f33075d4209e3300a06082a8648ce3d0403023010310e300c060355040a130553746f726a3022180f30303031303130313030303030305a180f30303031303130313030303030305a3010310e300c060355040a130553746f726a3059301306072a8648ce3d020106082a8648ce3d03010703420004e081e83ce2eb66625fa90bc50df5e9ec660bf58dd1bfb9f43184b6043b9fd62fc7ea64356cb778f0b54cbefb5c2cd6d89b7eb97b8f48f8f0f81e1944446764fca3383036300e0603551d0f0101ff04040302020430130603551d25040c300a06082b06010505070301300f0603551d130101ff040530030101ff300a06082a8648ce3d0403020348003045022100a245047d5fba2b7ba90b11ef47f477da0aede442592a5a33155e896c989eb498022016c17a5726c93ca2ef8b802733aa525664c10158b091a8cc4273cf9ac07ba1024a40d21be11d7f870387e57771f7c7341a3d6aeffaa5f9b67498f5bf4cc23352a89fa721f85da1c5ce25301fedd5f6d7fae7e57960fec2b8fc902948002846a27985109a051a20a55a5b2fdb1209011e0384051a2e86db45b1e2973d792c7d6152614836a0f80022e5023082016130820107a0030201020210722edac90a72a5a3b635f259ce461bb3300a06082a8648ce3d0403023010310e300c060355040a130553746f726a3022180f30303031303130313030303030305a180f30303031303130313030303030305a3010310e300c060355040a130553746f726a3059301306072a8648ce3d020106082a8648ce3d0301070342000429a7f2c871c0e440ff84eb579d6e5626632c79d8a8261b077e4bcec5db08203218a93a72b3ad53fa3e3ec9dabe1599c761600662a71974d3342db1bfbe6d1038a33f303d300e0603551d0f0101ff0404030205a0301d0603551d250416301406082b0601050507030106082b06010505070302300c0603551d130101ff04023000300a06082a8648ce3d0403020348003045022100bd33c392493fc00340f65f3fa5cf39c3798616b2f2c194de05db99edf11d69d602204b605acadf48bd81b651e13d6896a2ac890bb901f01e6b93be034bd587d35ab322dd023082015930820100a00302010202105c9de0e61af356ea0c62c883c039490f300a06082a8648ce3d0403023010310e300c060355040a130553746f726a3022180f30303031303130313030303030305a180f30303031303130313030303030305a3010310e300c060355040a130553746f726a3059301306072a8648ce3d020106082a8648ce3d03010703420004f8e126cef15cd3d1b5c2d6eccf5bd9c0b12fe05eec2722b375a5226a495b8b86fc9f74bc82dd8a4302a6f32f761550c13776f3896ffff5f462692b55e00127f6a3383036300e0603551d0f0101ff04040302020430130603551d25040c300a06082b06010505070301300f0603551d130101ff040530030101ff300a06082a8648ce3d040302034700304402206904cae9b20ff77fa11503836f5fb9f361155376dc9b7de9005b879133b8c4e702206078677e1fe4f033aefb046eaa43b7e7a77c975aefbeb65e1f8a6a04ad1775132a40e8d7bf6fa1558b0222c3f76adb8d5b99e23ca963600532fa0cb9a0ee49dbd8f7c260934fc07904fc34fc087d49c8c0491676231c0932db5b87796a6be672d9d2', '\xa55a5b2fdb1209011e0384051a2e86db45b1e2973d792c7d6152614836a0f800', 1, 666, '2019-02-14 10:09:54.420181-05', '2019-02-14 11:09:54-05');
INSERT INTO "irreparabledbs" ("segmentpath", "segmentdetail", "pieces_lost_count", "seg_damaged_unix_sec", "repair_attempt_count") VALUES ('\x49616d5365676d656e746b6579696e666f30', '\x49616d5365676d656e7464657461696c696e666f30', 10, 1550159554, 10);
INSERT INTO "injuredsegments" ("id", "info") VALUES (1, '\x0a0130120100');