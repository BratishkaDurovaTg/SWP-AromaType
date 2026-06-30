# Catalog Bot

Catalog management is handled by a separate Telegram bot process. The public web
app has no registration, JWT, or web admin panel.

## Access

The bot asks for a password after `/start`. The password is configured through
the production environment variable `CATALOG_BOT_PASSWORD`.

Do not commit the bot token or password to Git.

## Menu Flow

After password login, the bot shows menu buttons:

- `Каталог` - open the fragrance list and choose a product with inline buttons.
- `Добавить товар` - start the guided product creation flow.
- `Помощь` - show command reference.
- `Отмена` - cancel the current add/edit/photo flow.

Inside a product card, inline buttons allow editing fields, replacing the photo,
toggling `is_active`, and returning to the catalog.

## Commands

- `/list` - show all fragrances, including inactive items.
- `/view id` - show the full product card stored in the database.
- `/add` - create a new fragrance through a guided step-by-step flow.
- `/edit id` - show editable field buttons for a fragrance.
- `/set id field value` - update a single field.
- `/photo id` - upload or replace a product photo by sending a Telegram photo.
- `/toggle id` - switch product `is_active` on or off.
- `/cancel` - cancel the current flow.

## Editable Fields

- `name`
- `brand`
- `price`
- `volumes`
- `description`
- `top`
- `middle`
- `base`
- `accords`
- `psychotype`
- `scores`
- `active`
- `image_url`

## Examples

```text
/set miami-shake name Miami Shake
/set miami-shake price 8393
/set miami-shake volumes 50:8393, 100:12990
/set miami-shake top клубника, бергамот
/set miami-shake psychotype aesthetic
/set miami-shake scores drive:20, focus:35, aesthetic:90, power:25
/set miami-shake active yes
/photo miami-shake
```

After `/photo id`, send the image as a Telegram photo. The bot downloads it into
the shared `/uploads` volume and saves the public `/uploads/...` URL in the
product card.
