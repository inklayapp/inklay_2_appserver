---
description: Build your own Inklay dashboard on your computer, customized with Claude.
---

# Inklay App Builder

The Inklay App Server is a small app that runs on your computer and creates a dashboard image for your Inklay. Your Inklay shows it through the [Image URL](image-url.md) app.

You can change the dashboard however you like by describing it to Claude, for example *"show today's weather for Zurich and my next calendar events"*.

{% hint style="info" %}
The dashboard is created on your computer and sent directly to your Inklay over your local Wi-Fi. No data goes through the Inklay cloud.
{% endhint %}

## What you need

* An Inklay 2 with the Inklay app
* A Mac or Windows computer on the **same Wi-Fi** as your Inklay
* Optional, for customizing: [Claude Code](https://claude.com/claude-code) or the Claude desktop app

## Start the App Server

1. Download and unzip the **Inklay App Server** folder.
2. Start it:
   * **Mac:** double-click **Start Inklay App Server.command**.
   * **Windows:** double-click **dist/InklayAppServer-windows.exe**.
3. A browser page opens. It shows your **Image URL** and a preview of the dashboard.

A small window (Terminal on Mac) stays open while the App Server runs. Keep it open.

### First start on Mac

macOS may say the app "cannot be opened" because it was downloaded from the internet.

1. Right-click **Start Inklay App Server.command** and choose **Open**.
2. Click **Open** again.

If that doesn't work, go to **System Settings → Privacy & Security** and click **Open Anyway**.

If macOS asks whether the app may accept incoming network connections, click **Allow**. Otherwise your Inklay can't reach it.

### First start on Windows

1. If you see "Windows protected your PC", click **More info → Run anyway**.
2. When the firewall asks, allow access on **private networks**. Otherwise your Inklay can't reach it.

## Connect your Inklay

1. On the browser page, click the copy icon next to the **Image URL**. It looks like `http://192.168.0.33:8080/dashboard.png`.
2. In the Inklay app, add the [Image URL](image-url.md) app to a button.
3. Paste the URL and choose how often your Inklay should wake up.
4. Optionally, set the app as active.

Your Inklay now loads the dashboard every time it wakes up.

{% hint style="warning" %}
Your computer must be **on, awake and connected to the same Wi-Fi** when the Inklay wakes up. If it's asleep or the App Server is closed, the Inklay can't load a new image and keeps showing the last one.
{% endhint %}

## Customize your dashboard with Claude

1. Make a copy of the App Server folder, so you keep the original as a template.
2. Open the copy in Claude Code or the Claude desktop app.
3. Describe what you want to see, for example:
   * *"Show today's weather for Zurich."*
   * *"Show my next 5 calendar events."*
   * *"Show my solar production from Solar Manager."*
   * *"Make it portrait."*
4. Claude changes the dashboard, checks the result and rebuilds the app.
5. Start the App Server again to see your new dashboard.

The folder contains a `CLAUDE.md` file. It tells Claude about the Inklay display (1872 × 1404 pixels, greyscale), the default font and how to rebuild the app, so you don't have to explain it.

{% hint style="info" %}
If your dashboard needs a login or an API key, it is saved only on your computer, in `config.json`.
{% endhint %}

## Stop the App Server

Close its window, or click into it and press **Ctrl + C**. On a Mac, if asked, click **Terminate**.

## Troubleshooting

| Problem | Solution |
| --- | --- |
| The Inklay shows an old image | Check that the computer is awake, the App Server is running and both are on the same Wi-Fi. |
| The Inklay can't load the image | Allow incoming connections or private network access in the firewall (see first start above). |
| The URL changed | Your computer got a new IP address from the router. Copy the new URL into the Image URL app. Tip: reserve a fixed IP address for your computer in your router. |
| The image looks too dark or too bright | Adjust **Gamma** in the Image URL app settings. |
| The port is already in use | The App Server automatically uses the next free port (8081, 8082 …). Use the URL shown on the browser page. |
