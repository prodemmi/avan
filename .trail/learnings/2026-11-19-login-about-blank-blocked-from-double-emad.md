---
title: "Login about:blank#blocked from double-port form action"
author: Emad
date: 2026-11-19
branch: main
id: l-badee9be
scope: repo
status: active
src: agent
touches:
  - src/internal/handlers/auth.go
  - src/internal/repositories/repository.go
  - src/views/auth/login.templ
  - src/views/auth/login_templ.go
---

Login POST went to about:blank#blocked because AuthHandler.LoginPage built the form action with GetServiceURL(), which always appended ServerPort. docker-compose sets SERVER_DOMAIN=notification.layra.local:8083, so the action became http://notification.layra.local:8083:8083/auth/login, an invalid URL Chrome refuses to navigate to. Auth now uses origin-relative paths (/auth/login, /dashboard) and host-only cookies (empty Domain), since the app sits behind nginx and the cookie Domain attribute must not include scheme or port. GetServiceURL now detects an existing port instead of appending a second one.
