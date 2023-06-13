-- Insert board invites
INSERT INTO board_invites (id, receiver_id, board_id, sender_id, status, created_at, updated_at)
VALUES
    -- Board 1 invites
    ('cc143c5c-3a10-46b1-b734-8e2049b719ff', '3c85c78d-d6a2-48c1-aec0-7fe2a9e2d8db', 'b9e95ae4-9c3f-412f-8b3b-201bd7083fc1', 'd0865843-8494-4d6a-b9be-5c8f7d0e568f', 'ACCEPTED', '2023-05-29 10:30:00', '2023-05-29 10:30:00'),
    ('f9647cc4-4f2e-46e3-b5d9-345f79f6017e', '82d1a682-82e7-4ed9-bb6d-ef5b53f9827a', 'b9e95ae4-9c3f-412f-8b3b-201bd7083fc1', 'd0865843-8494-4d6a-b9be-5c8f7d0e568f','PENDING', '2023-05-29 10:30:00', '2023-05-29 10:30:00');