/*SELECT "evento"."evento_id","evento"."organizador_id","evento"."categoria_id","evento"."titulo","evento"."descripcion","evento"."lugar","evento"."evento_estado","evento"."cant_me_gusta","evento"."cant_no_interesa","evento"."cant_vendido_total","evento"."imagen_descripcion","evento"."imagen_portada","evento"."video_presentacion","evento"."imagen_escenario","evento"."total_recaudado","evento"."estado","evento"."usuario_creacion","evento"."fecha_creacion","evento"."usuario_modificacion","evento"."fecha_modificacion", f.fecha_evento FROM "evento" 
JOIN evento_fecha ef ON ef.evento_id = evento.evento_id 
JOIN fecha f ON f.fecha_id = ef.fecha_id 
WHERE f.fecha_evento >= CURRENT_DATE 
AND evento.evento_estado = 1 
AND evento.estado = 1 
AND f.fecha_evento = (
    SELECT MIN(f2.fecha_evento)
    FROM fecha f2
    JOIN evento_fecha ef2 ON ef2.fecha_id = f2.fecha_id
    WHERE ef2.evento_id = evento.evento_id
)
ORDER BY ((evento.cant_me_gusta - evento.cant_no_interesa) / GREATEST(1, (f.fecha_evento::date - CURRENT_DATE))) DESC;



SELECT evento_id, fecha.fecha_evento FROM evento_fecha JOIN fecha ON fecha.fecha_id = evento_fecha_id;


SELECT f.fecha_evento as fecha_, CURRENT_DATE as hoy, GREATEST(1, (f.fecha_evento::date - CURRENT_DATE)) as dif, ((evento.cant_me_gusta - evento.cant_no_interesa) / GREATEST(1, (f.fecha_evento::date - CURRENT_DATE))) AS score FROM fecha f
JOIN evento_fecha ef on ef.fecha_id = f.fecha_id
JOIN evento ON evento.evento_id = ef.evento_id 
WHERE ef.evento_id = 9
ORDER BY score DESC;

SELECT "evento"."evento_id","evento"."organizador_id","evento"."categoria_id","evento"."titulo","evento"."descripcion","evento"."lugar","evento"."evento_estado","evento"."cant_me_gusta","evento"."cant_no_interesa","evento"."cant_vendido_total","evento"."imagen_descripcion","evento"."imagen_portada","evento"."video_presentacion","evento"."imagen_escenario","evento"."total_recaudado","evento"."estado","evento"."usuario_creacion","evento"."fecha_creacion","evento"."usuario_modificacion","evento"."fecha_modificacion" FROM "evento" JOIN evento_fecha ef ON ef.evento_id = evento.evento_id JOIN fecha f ON f.fecha_id = ef.fecha_id WHERE f.fecha_evento >= CURRENT_DATE AND evento.evento_estado = 1 AND evento.estado = 1 AND f.fecha_evento = ( SELECT MIN(f2.fecha_evento)
                FROM fecha f2
                JOIN evento_fecha ef2 ON ef2.fecha_id = f2.fecha_id
                WHERE ef2.evento_id = evento.evento_id) ORDER BY ((2*evento.cant_me_gusta - evento.cant_no_interesa) / GREATEST(1, (f.fecha_evento::date - CURRENT_DATE))) DESC ;

SELECT * FROM interaccion;

SELECT * FROM usuario;

SELECT * FROM evento;

SELECT * FROM "evento" WHERE organizador_id = 1 AND evento_estado=1;

SELECT 
            COALESCE(SUM(DISTINCT oc.total), 0) AS ingreso_total,
            COALESCE(SUM(DISTINCT oc.monto_fee_servicio), 0) AS cargo_serv,
            COUNT(t.ticket_id) AS tickets_vendidos
         FROM orden_de_compra oc 
         JOIN ticket t ON t.orden_de_compra_id = oc.orden_de_compra_id 
         JOIN evento_fecha ef ON ef.evento_fecha_id = t.evento_fecha_id 
         WHERE ef.evento_id = 9 
                AND oc.estado_de_orden = 1 
                AND t.estado_de_ticket = 1 
                AND (oc.fecha BETWEEN '2025-11-28 00:00:00' AND '2025-11-28 23:59:00')
        GROUP BY oc.orden_de_compra_id, oc.total, oc.monto_fee_servicio;

SELECT * FROM orden_de_compra
WHERE usuario_id = 5;

SELECT * FROM evento_fecha
WHERE evento_fecha.evento_id = 9;

SELECT * FROM ticket
JOIN evento_fecha ef ON ef.evento_fecha_id = ticket.evento_fecha_id
WHERE ef.evento_id = 9;


SELECT s.sector_tipo AS sector,
        --s.total_entradas as capacidad,
        COUNT(t.ticket_id) AS tickets_vendidos,
        COALESCE(SUM(tf.precio), 0) AS ingresos
FROM ticket t 
JOIN orden_de_compra oc ON oc.orden_de_compra_id = t.orden_de_compra_id 
JOIN evento_fecha ef ON ef.evento_fecha_id = t.evento_fecha_id 
JOIN tarifa tf ON tf.tarifa_id = t.tarifa_id 
JOIN sector s ON s.sector_id = tf.sector_id 
WHERE ef.evento_id = 7 
AND oc.estado_de_orden = 1 
AND t.estado_de_ticket = 1 
AND (oc.fecha BETWEEN '2025-11-28 00:00:00' AND '2025-11-28 18:04:48.446') 
GROUP BY s.sector_id;--, s.total_entradas
*/

SELECT 
                        s.sector_tipo AS tipo_sector,
                        s.total_entradas as capacidad,
                        COUNT(t.ticket_id) AS tickets_vendidos,
                        COALESCE(SUM(tf.precio), 0) AS ingresos
                 FROM sector s JOIN tarifa tf ON tf.sector_id = s.sector_id JOIN ticket t ON t.tarifa_id = tf.tarifa_id JOIN orden_de_compra oc ON oc.orden_de_compra_id = t.orden_de_compra_id WHERE s.evento_id = 4 AND oc.estado_de_orden = 1 AND t.estado_de_ticket = 1 AND (DATE(oc.fecha) BETWEEN '2025-11-28 00:00:00' AND '2025-11-28 19:00:00') GROUP BY s.sector_tipo, s.total_entradas;

SELECT * FROM orden_de_compra;
