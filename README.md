NestGo (NG) is a lightweight framework inspired by the NestJS controller pattern, designed to provide a structured, declarative way to build HTTP applications while remaining framework-agnostic.

Key Features
🧩 Decorator-Based Architecture

Uses the decorator pattern to define controllers, routes, and behaviors

Promotes clean separation of concerns and readable code structure

🔌 Framework Agnostic

Easily switch between HTTP frameworks such as Fiber, Echo, Gin or others

Core logic is decoupled from the underlying HTTP implementation

🔁 Predictable Request Lifecycle

Requests flow through a well-defined pipeline:

Request
→ Middleware
→ Guard
→ Interceptor
→ Handler

This mirrors NestJS’s execution model while remaining simple and extensible.

🧠 Metadata Extraction

Extract and manage:

Application metadata

Controller metadata

Route metadata

Enables advanced use cases like:

Automatic routing

API documentation

Middleware and interceptor orchestration

Design Goals

Familiar developer experience for NestJS users

Minimal core with extensibility in mind

Clear execution flow and strong abstractions

Portable across HTTP frameworks
