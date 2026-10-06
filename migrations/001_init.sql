CREATE DATABASE IF NOT EXISTS job_tracker
    CHARACTER SET utf8mb4
    COLLATE utf8mb4_unicode_ci;

USE job_tracker;

CREATE TABLE IF NOT EXISTS applications (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT , 
    company VARCHAR(100) NOT NULL ,
    position VARCHAR(100) NOT NULL ,
    status VARCHAR(20) NOT NULL , 
    applied_at DATE NULL ,
    notes VARCHAR(1000) NOT NULL DEFAULT '' , 
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
        ON UPDATE CURRENT_TIMESTAMP,

    PRIMARY KEY (id) ,
    INDEX idx_applications_status (status),
    INDEX idx_applications_company (company),
    INDEX idx_applications_applied_at (applied_at)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;