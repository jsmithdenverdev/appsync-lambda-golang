# appsync-lambda-golang

A sample AWS AppSync GraphQL API leveraging lambda functions as resolvers.

## GraphQL Schema

The resulting GraphQL API has the following schema.

```graphql
input CreateItemInput {
  name: String!
}

type Item {
  name: String!
}

type CreateItemPayload {
  item: Item!
}

type Query {
  item(id: ID!): Item
  items: [Item!]!
}

type Mutation {
  createItem(input: CreateItemInput!): CreateItemPayload!
}
```

## Resolvers

| type     | field        | lambda                                    | mode         |
| -------- | ------------ | ----------------------------------------- | ------------ |
| Mutation | `createItem` | [CreateItemFunction](#createitemfunction) | invoke       |
| Query    | `item`       | [ReadItemFunction](#readitemfunction)     | invoke       |
| Query    | `items`      | [ListItemsFunction](#listitemsfunction)   | batch invoke |

## Functions

### CreateItemFunction

A lambda function serving the `createItem` mutation. Returns `responses.CreateItem` using the `Name` from `requests.CreateItem`.

| meta    | value                                                                |
| ------- | -------------------------------------------------------------------- |
| runtime | go1.x                                                                |
| uri     | [cmd/create-item/main.go](cmd/create-item/main.go)                   |
| handler | [internal/handlers/create_item.go](internal/handlers/create_item.go) |

### ReadItemFunction

A lambda function serving the `item` query. Returns `responses.Item`.

#### ⛔️ Note: Currently returns a not implemented error.

| meta    | value                                                            |
| ------- | ---------------------------------------------------------------- |
| runtime | go1.x                                                            |
| uri     | [cmd/read-item/main.go](cmd/read-item/main.go)                   |
| handler | [internal/handlers/read_item.go](internal/handlers/read_item.go) |

### ListItemsFunction

A lambda function serving the `items` query. Returns `responses.Item[]`. This handler is invoked as a batch invoke and operates in a similar fashion to a data loader.

#### ⛔️ Note: Currently returns a not implemented error.

| meta    | value                                                              |
| ------- | ------------------------------------------------------------------ |
| runtime | go1.x                                                              |
| uri     | [cmd/list-items/main.go](cmd/list-items/main.go)                   |
| handler | [internal/handlers/list_items.go](internal/handlers/list_items.go) |

## Building

`sam build`

## Deploying

Initial: `sam deploy --guided`

Subsequent: `sam deploy`
