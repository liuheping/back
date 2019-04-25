-- +goose Up
CREATE TABLE brand_series (
  id VARCHAR(32) NOT NULL
  constraint brand_series_pkey
			primary key,
  brand_id VARCHAR(32) NOT NULL,
  series VARCHAR(64) NOT NULL,
  image VARCHAR(255),
  machine_types VARCHAR[] NOT NULL,
  sort_order INT NOT NULL DEFAULT 255
);
-- bcsvpsdm54kmvd4pec4g	SK250LC	brand/2018-07-02/bct07odm54kncsf31d3g.png	{SK250LC-8}

-- +goose Down
DROP TABLE brand_series;
