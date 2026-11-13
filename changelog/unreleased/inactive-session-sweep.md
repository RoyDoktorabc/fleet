---
type: improved
---

- **Old sessions park themselves.** A session that sits idle, finished or errored for 48 hours is now suspended without waiting for memory pressure. Press `Enter` to resume it. Set `"suspend_inactive_after_hours"` in `~/.config/fleet/config.json` to change the limit, or `0` to turn it off.
- **Dead rows look dead.** A session whose agent exited shows a red `✕` with a dimmed title, instead of a blue dot that read as finished.
