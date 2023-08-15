vmangos realmd process

1. connect to db
2. prefetch the realm list
3. update bans
4. start TCP server on port 3724
5. accept connections
6. wait for client to send packet
7. read first byte of packet to determine command

login challenge (opcode 0x00)

1. read packet, convert endian, insert into struct
2. reject login if IP banned
3. reject login if email not verified
4. verify SRP6 v/s values match database
5. respond with some SRP6 data
6. wait for login proof

login proof (opcode 0x01)

1. check if build matches requested version
2. perform SRP6 proof check
3. if success, update account in db with session, set client as authed
4. if fail, temp ban user if too many login attempts, then reject
