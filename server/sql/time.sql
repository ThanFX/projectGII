SELECT real_time, world_time, id, time_speed
  FROM public."time";

/* Количество секунд мира на момент старта новой разработки */
UPDATE time SET world_time = 439555701 WHERE id = 1;

SELECT * FROM config;
SELECT * FROM persons;
SELECT * FROM person_health_characteristic;

INSERT INTO public.person_health_characteristic(
            person_id, health, fatigue, hunger, thirst, somnolency)
    VALUES (1, 100.0, 0.0, 0.0, 0.0, 0.0);

UPDATE public.config SET 
value = '{"State": 10, "HTS": 10, "ED": 10}'
WHERE id = 'check_periods';