# Admin Workflows & Examples

Practical examples untuk menggunakan fitur admin bot café.

## 🔐 Admin Authentication

### Setup Admin Access

**1. Get your Telegram ID:**
- Chat dengan @userinfobot di Telegram
- Bot akan kirim ID Anda (contoh: `647143027`)

**2. Add to `.vars.json`:**
```json
{
  "admin_telegram_ids": ["647143027"],
  "admin_usernames": ["your_telegram_username"]
}
```

**3. Restart bot:**
```bash
make stop
make dev-local-hot
```

**4. Test:**
- Open Telegram
- Send `/start` to your bot
- Should see "Anda login sebagai *Admin*"
- Admin panel buttons akan muncul

## 📋 Menu Management

### Add New Menu Item

**Step-by-step:**

1. **Start:** `/admin` atau klik "👨‍💼 Panel Admin"
2. **Click:** "📋 Kelola Menu"
3. **Click:** "➕ Tambah Menu Baru"
4. **Enter name:** `Cappuccino`
5. **Enter price:** `25000`
6. **Choose category:** `Coffee`
7. **Enter description:** `Kopi susu premium dengan busa lembut`
8. **Success:** Bot shows confirmation with menu details

**Complete chat example:**
```
Bot: Masukkan nama menu:
You: Cappuccino

Bot: Masukkan harga menu (angka saja, tanpa Rp):
You: 25000

Bot: Pilih kategori:
• Coffee
• Tea
• Snacks

Ketik nama kategori:
You: Coffee

Bot: Masukkan deskripsi menu (atau ketik - untuk skip):
You: Kopi susu premium dengan busa lembut

Bot: ✅ Menu berhasil ditambahkan!

🍽️ Cappuccino
💰 Rp 25.000
```

### View All Menus

1. Click "📋 Kelola Menu"
2. Click "📖 Lihat Semua Menu"
3. Bot shows list with:
   - ✅/❌ Status (available/not available)
   - Name
   - Price
   - Category

### Delete Menu

1. Click "📋 Kelola Menu"
2. Click "🗑️ Hapus Menu"
3. Bot shows menu list with delete buttons
4. Click "🗑️ Hapus [Menu Name]"
5. Confirm: "✅ Ya, Hapus"
6. Success confirmation

**Cancel anytime:** Send `/cancel`

## 🎉 Promo Management

### Create New Promo

**Example: 20% Discount Promo**

```
Bot: Masukkan judul promo:
You: Diskon Weekend 20%

Bot: Masukkan deskripsi promo (atau - untuk skip):
You: Diskon 20% untuk semua menu setiap weekend

Bot: Pilih tipe diskon:
Ketik: percentage atau amount
You: percentage

Bot: Masukkan jumlah diskon (dalam %, angka saja):
You: 20

Bot: Masukkan tanggal mulai (format: YYYY-MM-DD):
You: 2025-12-01

Bot: Masukkan tanggal akhir (format: YYYY-MM-DD):
You: 2025-12-31

Bot: ✅ Promo berhasil ditambahkan!

🎁 Diskon Weekend 20%
```

**Example: Fixed Amount Discount**

```
Tipe diskon: amount
Jumlah diskon: 10000
→ Result: Diskon Rp 10.000
```

### View Active Promos

1. Click "🎉 Kelola Promo"
2. Click "📖 Lihat Semua Promo"
3. Shows all promos with:
   - ✅ Active / ❌ Inactive
   - Title
   - Discount value

### Delete Promo

Similar to delete menu:
1. "🎉 Kelola Promo" → "🗑️ Hapus Promo"
2. Select promo to delete
3. Confirm deletion

## ℹ️ Café Info Management

### View Current Info

1. Click "ℹ️ Kelola Info Café"
2. Click "📖 Lihat Info Café"
3. Shows:
   - Nama café
   - Alamat
   - Telepon
   - Jam operasional
   - Deskripsi

### Edit Café Info

**Example: Update Cafe Name**

```
1. Click "ℹ️ Kelola Info Café"
2. Click "✏️ Edit Info Café"
3. Bot shows current info with field buttons
4. Click "📍 Edit Nama"

Bot: Masukkan nama café baru:
(Ketik /cancel untuk membatalkan)

You: Café Nusantara Premium

Bot: ✅ Info Café berhasil diperbarui!

📍 Nama: Café Nusantara Premium
🏠 Alamat: (unchanged)
📞 Telepon: (unchanged)
...
```

**Editable fields:**
- 📍 Nama
- 🏠 Alamat
- 📞 Telepon
- 📧 Email
- 🕐 Jam Buka
- 🕔 Jam Tutup
- 📝 Deskripsi

**Update operating hours:**
```
Jam Buka: 08:00
Jam Tutup: 22:00
```

**Clear optional fields:**
Type `-` to clear email or description.

## 📁 Category Management

### Add New Category

```
1. Click "📁 Kelola Kategori"
2. Click "➕ Tambah Kategori Baru"

Bot: Masukkan nama kategori:
You: Desserts

Bot: ✅ Kategori Desserts berhasil ditambahkan!
```

### View Categories

1. Click "📁 Kelola Kategori"
2. Click "📖 Lihat Semua Kategori"
3. Shows list of all categories

### Delete Category

⚠️ **Warning:** Deleting a category may affect menus using it!

```
1. Click "📁 Kelola Kategori"
2. Click "🗑️ Hapus Kategori"
3. Select category to delete
4. Confirm deletion
```

## 🔄 Common Workflows

### Daily Operations

**Morning - Setup promos:**
```
1. /admin
2. Kelola Promo → Tambah Promo
3. Create "Early Bird Discount"
4. Set valid dates
```

**Update menu items:**
```
1. Kelola Menu → Lihat Semua Menu
2. Check availability
3. Disable sold-out items (when edit feature available)
```

**Evening - Review:**
```
1. Check active promos
2. Update info if needed
3. Plan tomorrow's offers
```

### Weekly Tasks

**Monday:**
- Review and update weekly promos
- Check menu popularity
- Update café info if needed

**Friday:**
- Create weekend special promos
- Update operating hours if special schedule

### Monthly Tasks

**Monthly review:**
- Clean up expired promos
- Update menu prices
- Review category organization
- Update café description for special events

## 💡 Pro Tips

### 1. Use Descriptive Names

**Good:**
```
Menu: "Cappuccino (Medium)"
Promo: "Weekend Lunch Discount 15%"
Category: "Hot Beverages"
```

**Avoid:**
```
Menu: "Coffee 1"
Promo: "Promo A"
Category: "Cat1"
```

### 2. Plan Promo Dates

```
# Holiday promo
Start: 2025-12-24
End: 2025-12-26

# Month-long promo
Start: 2025-12-01
End: 2025-12-31
```

### 3. Organize Categories

Create logical groupings:
- Hot Beverages
- Cold Beverages
- Snacks
- Main Course
- Desserts

### 4. Keep Info Updated

Update café info when:
- Hours change
- Contact changes
- Special events
- Seasonal adjustments

### 5. Use Cancel Command

Anytime during dialog:
```
/cancel
```
Cancels current operation and clears state.

## 🐛 Troubleshooting

### Can't Access Admin Panel

**Check:**
1. Telegram ID in `.vars.json`
2. Bot restarted after config change
3. Using `/start` command

**Fix:**
```bash
cat .vars.json  # Verify ID
make stop
make dev-local-hot
```

### Promo Not Showing for Users

**Causes:**
- End date passed
- Not marked as active
- Date format wrong

**Fix:**
- Check dates: YYYY-MM-DD format
- Ensure `is_active: true`
- Create new promo if needed

### Menu Changes Not Reflected

**Solution:**
```bash
# Restart bot
make stop
make dev-local-hot

# Or with Docker
docker restart cafe-bot-agent
```

## 📊 Admin Best Practices

### 1. Backup Before Major Changes

```bash
# Backup databases
tar -czf backup-$(date +%Y%m%d).tar.gz data/
```

### 2. Test in Staging First

If you have staging bot:
- Test new promos
- Verify menu changes
- Check info updates

### 3. Document Changes

Keep track of:
- When promos were created
- Menu price changes
- Category updates

### 4. Regular Cleanup

Monthly:
- Delete expired promos
- Remove discontinued menus
- Update categories

### 5. Multiple Admins

Add multiple admin IDs in `.vars.json`:
```json
{
  "admin_telegram_ids": [
    "123456789",
    "987654321"
  ]
}
```

## 📱 Mobile vs Desktop

**Mobile (Primary):**
- Fast menu management
- Quick promo creation
- On-the-go updates

**Desktop (Planning):**
- Easier for bulk data entry
- Better for long descriptions
- Good for planning campaigns

Both work the same way - use what's convenient!

---

**Happy Managing!** 🎉

These workflows make managing your café bot efficient and organized. Practice a few times and it becomes second nature!

For technical issues, see [Troubleshooting Guide](../guides/troubleshooting.md).
