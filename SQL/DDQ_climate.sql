/******************************************************************************
    CLIMATEFIND DATA DEFINITION QUERIES
******************************************************************************/
DROP TABLE IF EXISTS 
users_have_activities, 
users_have_answers, 
answers_have_questions,
users_have_goals,
surveys_have_questions,
users,
activities,
answers,
surveys,
leaderboards,
leaderboards_have_users,
questions,
goals;

/******************************************************************************
    USERS
    Note that user is inserted without karma to let default of 0 take effect
******************************************************************************/
CREATE TABLE IF NOT EXISTS `users` (
	`id` int NOT NULL AUTO_INCREMENT,
	`account` varchar(255) NOT NULL,
	`pass` varchar(255),
	`first_name` varchar(255) NOT NULL,
	`last_name` varchar(255),
    `karma` int DEFAULT 0, 
    PRIMARY KEY(`id`),
	UNIQUE(`account`)
	) ENGINE=InnoDB;

INSERT INTO users (account, pass, first_name, last_name)
VALUES ("test one", "test_pass 1", "one", "lastname 1");

INSERT INTO users (account, pass, first_name, last_name)
VALUES ("test two", "test_pass 2", "two", "lastname 2");

INSERT INTO users (account, pass, first_name, last_name)
VALUES ("test three", "test_pass 3", "three", "lastname 3");

INSERT INTO users (account, pass, first_name, last_name)
VALUES ("test 4", "test_pass 4", "four", "lastname 4");

/******************************************************************************
    GOALS
******************************************************************************/

CREATE TABLE IF NOT EXISTS `goals` (
	`id` int NOT NULL AUTO_INCREMENT,
	`name` varchar(255) NOT NULL,
    PRIMARY KEY(`id`),
	UNIQUE(`name`)
	) ENGINE=InnoDB;

INSERT INTO goals (name)
VALUES ("goal one");

INSERT INTO goals (name)
VALUES ("goal two");

/******************************************************************************
    SURVEYS
******************************************************************************/

CREATE TABLE IF NOT EXISTS `surveys` (
	`id` int NOT NULL AUTO_INCREMENT,
	`name` varchar(255) NOT NULL,
    PRIMARY KEY(`id`),
	UNIQUE(`name`)
	) ENGINE=InnoDB;

INSERT INTO surveys (name)
VALUES ("some survey");

/******************************************************************************
    QUESTIONS
******************************************************************************/

CREATE TABLE IF NOT EXISTS `questions` (
	`id` int NOT NULL AUTO_INCREMENT,
	`goal_id` int NOT NULL,
    `question_text` varchar(255) NOT NULL,
    PRIMARY KEY(`id`),
	FOREIGN KEY (`goal_id`) REFERENCES goals(`id`)
	) ENGINE=InnoDB;

INSERT INTO questions (goal_id, question_text)
VALUES (1, "question text here");

/******************************************************************************
    ANSWERS
******************************************************************************/

CREATE TABLE IF NOT EXISTS `answers` (
	`id` int NOT NULL AUTO_INCREMENT,
	`question_id` int NOT NULL,
    `answer_text` varchar(255) NOT NULL,
    `answer_code` int NOT NULL,
    PRIMARY KEY(`id`),
	FOREIGN KEY (`question_id`) REFERENCES questions(`id`)
	) ENGINE=InnoDB;

INSERT INTO answers (question_id, answer_text, answer_code)
VALUES (1, "answer text here", 5);

/******************************************************************************
    ACTIVITIES
******************************************************************************/

CREATE TABLE IF NOT EXISTS `activities` (
	`id` int NOT NULL AUTO_INCREMENT,
	`goal_id` int NOT NULL,
    `activity` varchar(255) NOT NULL,
    `karma_value` int NOT NULL,
    PRIMARY KEY(`id`),
	FOREIGN KEY (`goal_id`) REFERENCES goals(`id`)
	) ENGINE=InnoDB;

INSERT INTO activities (goal_id, activity, karma_value)
VALUES (1, "activity name", 5);

/******************************************************************************
    LEADERBOARDS
******************************************************************************/

CREATE TABLE IF NOT EXISTS `leaderboards` (
	`id` int NOT NULL AUTO_INCREMENT,
    `title` varchar(255) NOT NULL,
    `goal_id` int NOT NULL,
    PRIMARY KEY(`id`),
    UNIQUE(`title`),
    FOREIGN KEY(`goal_id`) REFERENCES goals(`id`)
	) ENGINE=InnoDB;

INSERT INTO leaderboards (title, goal_id)
VALUES ("first board", 1);
    
INSERT INTO leaderboards (title, goal_id)
VALUES ("second board", 2);

/******************************************************************************
    LEADERBOARDS_HAVE_USERS
******************************************************************************/
/*CREATE TABLE leaderboards_have_users AS 
SELECT leaderboards.id as lb_ID, leaderboards.title as leaderboard, leaders.id as user_id, leaders.karma FROM leaderboards
JOIN (SELECT users.id, users.karma, users_karma_goals.goal_id FROM users
JOIN (SELECT id, karma, users_g1.goal_id FROM users
JOIN (SELECT ug.user_id, ug.goal_id FROM users_have_goals ug
JOIN leaderboards lb
ON ug.goal_id=lb.goal_id) as users_g1
ON users_g1.user_id=users.id) as users_karma_goals
ON users.id = users_karma_goals.id) as leaders
ON leaders.goal_id = leaderboards.goal_id;
*/
/*
DROP TABLE IF EXISTS experiment;

CREATE TABLE experiment AS SELECT users.id, users.karma FROM users
JOIN (SELECT id, karma, users_g1.goal_id FROM users
JOIN (SELECT ug.user_id, ug.goal_id FROM users_have_goals ug
JOIN leaderboards lb
ON ug.goal_id=lb.goal_id) as users_g1
ON users_g1.user_id=users.id) as users_karma_goals
ON users.id = users_karma_goals.id AND users_karma_goals.goal_id = 1
ORDER BY users.karma desc;
*/
/**** 
This query just makes a leaderboard for goal #1 - subsitute goal number as needed
*****/
/*
SELECT users.id, users.karma, users_karma_goals.goal_id FROM users
JOIN (SELECT id, karma, users_g1.goal_id FROM users
JOIN (SELECT ug.user_id, ug.goal_id FROM users_have_goals ug
JOIN leaderboards lb
ON ug.goal_id=lb.goal_id) as users_g1
ON users_g1.user_id=users.id) as users_karma_goals
ON users.id = users_karma_goals.id AND users_karma_goals.goal_id = 1
ORDER BY users.karma desc;
*/
/******************************************************************************
    USERS_HAVE_GOALS
******************************************************************************/

CREATE TABLE IF NOT EXISTS `users_have_goals` (
	`id` int NOT NULL AUTO_INCREMENT,
	`goal_id` int NOT NULL,
    `user_id` int NOT NULL,
    PRIMARY KEY(`id`),
	FOREIGN KEY (`goal_id`) REFERENCES goals(`id`),
    FOREIGN KEY (`user_id`) REFERENCES users(`id`)
	) ENGINE=InnoDB;

INSERT INTO users_have_goals (goal_id, user_id)
VALUES (1, 1);
INSERT INTO users_have_goals (goal_id, user_id)
VALUES (1, 2);
INSERT INTO users_have_goals (goal_id, user_id)
VALUES (1, 3);
INSERT INTO users_have_goals (goal_id, user_id)
VALUES (1, 4);
/******************************************************************************
    SURVEYS_HAVE_QUESTIONS
******************************************************************************/

CREATE TABLE IF NOT EXISTS `surveys_have_questions` (
	`id` int NOT NULL AUTO_INCREMENT,
	`survey_id` int NOT NULL,
    `question_id` int NOT NULL,
    PRIMARY KEY(`id`),
	FOREIGN KEY (`survey_id`) REFERENCES surveys(`id`),
    FOREIGN KEY (`question_id`) REFERENCES questions(`id`)
	) ENGINE=InnoDB;

INSERT INTO surveys_have_questions (survey_id, question_id)
VALUES (1, 1);

/******************************************************************************
    USERS_HAVE_ANSWERS
******************************************************************************/

CREATE TABLE IF NOT EXISTS `users_have_answers` (
	`id` int NOT NULL AUTO_INCREMENT,
	`user_id` int NOT NULL,
    `answer_id` int NOT NULL,
    PRIMARY KEY(`id`),
	FOREIGN KEY (`user_id`) REFERENCES users(`id`),
    FOREIGN KEY (`answer_id`) REFERENCES answers(`id`)
	) ENGINE=InnoDB;

INSERT INTO users_have_answers (user_id, answer_id)
VALUES (1, 1);

/******************************************************************************
    USERS_HAVE_ACTIVITIES
******************************************************************************/

CREATE TABLE IF NOT EXISTS `users_have_activities` (
	`id` int NOT NULL AUTO_INCREMENT,
	`user_id` int NOT NULL,
    `activity_id` int NOT NULL,
    PRIMARY KEY(`id`),
	FOREIGN KEY (`user_id`) REFERENCES users(`id`),
    FOREIGN KEY (`activity_id`) REFERENCES activities(`id`)
	) ENGINE=InnoDB;

INSERT INTO users_have_activities (user_id, activity_id)
VALUES (1, 1);