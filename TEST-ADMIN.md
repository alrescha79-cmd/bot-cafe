# 🧪 Test Admin Authentication

## Problem
Admin dengan ID `647143027` dan username `Alrescha79` tidak mendapat akses admin panel.

## Debugging Steps

### 1. Start Services (Terminal 1)
```bash
cd /home/son/Projects/bot-cafe
make dev-local
```

### 2. Start Agent with Logging (Terminal 2)
```bash
cd /home/son/Projects/bot-cafe
./scripts/test-agent-only.sh
```

### 3. Test di Telegram
Kirim: `/start`

### 4. Check Log Output
Di Terminal 2, akan muncul log seperti:
```
INFO: Checking admin: userID=647143027, username=Alrescha79
INFO: Admin IDs in config: [647143027]
INFO: Admin Usernames in config: [Alrescha79]
INFO: User 647143027 matched by ID
INFO: User 647143027 recognized as ADMIN, showing admin panel
```

## Expected Behavior
- ✅ User ID harus match dengan `.vars.json`
- ✅ Bot harus tampilkan tombol "👨‍💼 Panel Admin"
- ✅ Welcome message: "Anda login sebagai *Admin*"

## Jika Masih Gagal
Cek hal berikut:
1. Username di Telegram sama dengan `.vars.json` (case-sensitive!)
2. User ID benar (bisa cek dengan bot @userinfobot)
3. File `.vars.json` di-load dengan benar (cek log saat startup)

## Quick Fix
Jika username tidak match, gunakan User ID saja:
```json
{
  "admin_telegram_ids": ["647143027"],
  "admin_usernames": []
}
```
