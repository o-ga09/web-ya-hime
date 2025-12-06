GRANT ALL PRIVILEGES ON test_web_ya_hime.* TO 'user'@'%';
GRANT ALL PRIVILEGES ON develop_web_ya_hime.* TO 'user'@'%';
FLUSH PRIVILEGES;

CREATE DATABASE IF NOT EXISTS test_web_ya_hime;
CREATE DATABASE IF NOT EXISTS develop_web_ya_hime;
