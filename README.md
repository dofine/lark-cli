# lark

A CLI tool for interacting with Lark/Feishu APIs, designed for use with Claude Code and other AI assistants.

## Features

- **Calendar** - List, create, update, delete events; check availability; find common free time; RSVP
- **Contacts** - Look up users by ID, search by name, list department members
- **Documents** - Read documents as markdown, list folders, resolve wiki nodes, get comments
- **Messages** - Retrieve chat history, download attachments, send messages, add/list/remove reactions
- **Mail** - Read and search emails via IMAP with local caching
- **Minutes** - Get meeting recording metadata, export transcripts, download media

## Fork Differences

This is a fork of [yjwong/lark-cli](https://github.com/yjwong/lark-cli) with the following changes:

### Feishu (China) Region Support

The original project only supports Lark (international). This fork adds support for Feishu (China) by setting `region: "feishu"` in config or `LARK_REGION=feishu` as environment variable. This switches all OAuth and API endpoints from `larksuite.com` to `feishu.cn`.

### Reduced OAuth Scopes

Some scopes require company admin approval on Feishu, which blocks normal login. This fork removes those scopes so you can authenticate without waiting for admin approval. The trade-off is that the following features are **unavailable**:

| Removed Scope | Impact |
|---|---|
| `contact:department.organize:readonly` | Cannot view department org structure |
| `docs:doc:readonly` | Cannot read legacy Docs (new Docx still works) |
| `drive:drive:readonly` | Cannot browse Drive root folder |
| `space:document:retrieve` | Cannot retrieve wiki space documents |
| `mail:*` (4 scopes) | Mail commands do not work |
| `minutes:*` (2 scopes) | Minutes commands do not work |

### Additional Features

For a full list of enhancements added in this fork beyond the above, see the [GitHub Issues](https://github.com/dofine/lark-cli/issues?q=is%3Aissue).

## Why This Tool?

The official Lark MCP server exists, but its tools are not token-efficient. Each tool call returns verbose responses that consume significant context window space when used with AI assistants.

This CLI addresses that by:

- **Returning compact JSON** - Structured output optimized for programmatic consumption
- **Providing markdown conversion** - Documents are converted to markdown (~2-3x smaller than raw block structures)
- **Supporting selective queries** - Fetch only what you need (e.g., just event IDs, just document titles)

The result: AI assistants can interact with Lark using fewer tokens, leaving more context for actual work.

## Quick Start

1. Create an app at the Developer Console:
   - Lark: https://open.larksuite.com
   - Feishu: https://open.feishu.cn
2. Copy `config.example.yaml` to `.lark/config.yaml` and add your App ID
3. Set `LARK_APP_SECRET` environment variable
4. (Feishu users) Set `LARK_REGION=feishu` or add `region: "feishu"` to config
5. Run `./lark auth login` to authenticate
6. Start using: `./lark cal list --week`

See [USAGE.md](USAGE.md) for full documentation.

## Building

```bash
make build    # Build binary to ./lark
make test     # Run tests
make install  # Install to $GOPATH/bin
```

## Usage with Claude Code

This tool is designed to be invoked via Claude Code skills. Pre-built skill definitions are included in the `skills/` directory.

### Installing Skills

Copy the skill directories to your Claude Code skills location:

```bash
# Project-specific (recommended)
cp -r skills/* /path/to/your/project/.claude/skills/

# Or user-wide
cp -r skills/* ~/.claude/skills/
```

Available skills:
- `calendar` - Manage calendar events, check availability, RSVP
- `contacts` - Look up users and departments
- `documents` - Read documents, list folders, browse wikis
- `messages` - Retrieve chat history, download attachments, send messages to users and chats
- `email` - Read and search emails via IMAP with local caching
- `minutes` - Get meeting recordings, export transcripts, download media

### Configuration

The skills assume `lark` is in your PATH. If not, you can either:

1. Add the binary location to your PATH
2. Edit the skill files to use the full path
3. Set `LARK_CONFIG_DIR` environment variable to point to your `.lark/` config directory

The JSON output format makes it straightforward for AI assistants to parse responses and take action.

## License

MIT - see [LICENSE](LICENSE)
