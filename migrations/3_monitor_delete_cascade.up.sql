-- Add cascading deletes for monitor-related records to allow monitor removal.
ALTER TABLE monitor_notifications DROP CONSTRAINT fk_monitor_notifications_monitor_id_monitors_id;
ALTER TABLE monitor_notifications
    ADD CONSTRAINT fk_monitor_notifications_monitor_id_monitors_id
    FOREIGN KEY (monitor_id) REFERENCES monitors(id) ON DELETE CASCADE;

ALTER TABLE pings DROP CONSTRAINT fk_pings_monitor_id_monitors_id;
ALTER TABLE pings
    ADD CONSTRAINT fk_pings_monitor_id_monitors_id
    FOREIGN KEY (monitor_id) REFERENCES monitors(id) ON DELETE CASCADE;

ALTER TABLE incidents DROP CONSTRAINT fk_incidents_monitor_id_monitors_id;
ALTER TABLE incidents
    ADD CONSTRAINT fk_incidents_monitor_id_monitors_id
    FOREIGN KEY (monitor_id) REFERENCES monitors(id) ON DELETE CASCADE;

ALTER TABLE incident_events DROP CONSTRAINT fk_incident_events_incident_id_incidents_id;
ALTER TABLE incident_events
    ADD CONSTRAINT fk_incident_events_incident_id_incidents_id
    FOREIGN KEY (incident_id) REFERENCES incidents(id) ON DELETE CASCADE;
