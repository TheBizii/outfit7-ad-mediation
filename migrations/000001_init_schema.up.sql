CREATE TABLE priority_lists (
	id SERIAL PRIMARY KEY,
	country_code VARCHAR(3) NOT NULL,
	ad_type VARCHAR(20) NOT NULL CHECK (ad_type IN ('banner', 'interstitial', 'rewarded_video')),
	last_updated TIMESTAMP NOT NULL DEFAULT NOW(),
	UNIQUE (country_code, ad_type)
);

CREATE TABLE priority_networks (
	id SERIAL PRIMARY KEY,
	priority_list_id INT NOT NULL REFERENCES priority_lists(id) ON DELETE CASCADE,
	network_name VARCHAR(50) NOT NULL,
	score REAL NOT NULL,
	UNIQUE (priority_list_id, network_name)
);
