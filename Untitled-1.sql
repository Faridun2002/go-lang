USE `astet`;

CREATE TABLE `lec` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name_lec` text DEFAULT NULL,
  `proiz_lec` text DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8;

insert  into `lec`(`id`,`name_lec`,`proiz_lec`) values 

(1,'Гипчалгин','Таджикистан'),
(2,'Парацетамол','Россия'),
(3,'dvsvsdv','sdvsdvsdv'),
(4,'jdgjg','jdghjdj');

