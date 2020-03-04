START TRANSACTION;

CREATE TABLE admin_studio (
  user_id VARCHAR(40) NOT NULL,
  studio_id VARCHAR(40) NOT NULL,
  PRIMARY KEY (user_id, studio_id),
  CONSTRAINT fk_studio_id_of_admin FOREIGN KEY (studio_id)
  REFERENCES studios(id)
);

CREATE TABLE invitations (
  email VARCHAR(255) NOT NULL,
  group_id VARCHAR(40) NOT NULL,
  is_accepted BOOLEAN DEFAULT false NOT NULL,
  group_type ENUM('studio', 'song') NOT NULL,
  PRIMARY KEY (email, group_id)
);

CREATE TABLE member_studio (
  user_id VARCHAR(40) NOT NULL,
  studio_id VARCHAR(40) NOT NULL,
  PRIMARY KEY (user_id, studio_id),
  CONSTRAINT fk_studio_id_of_member FOREIGN KEY (studio_id)
  REFERENCES studios(id)
);

CREATE TABLE guest_song (
  user_id VARCHAR(40) NOT NULL,
  song_id VARCHAR(40) NOT NULL,
  PRIMARY KEY (user_id, song_id),
  CONSTRAINT fk_song_id_of_guest FOREIGN KEY (song_id)
  REFERENCES songs(id)
);

COMMIT;
