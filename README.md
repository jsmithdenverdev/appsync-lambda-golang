# appsync-lambda-golang

A sample AWS AppSync GraphQL API leveraging lambda functions as resolvers.

## Why?

GraphQL (when done right) is a fantastic way to build highly expressive API's. Having the flexibility to back resolvers using lambda functions unlocks a new level of customization and scalability. Using this approach we can independently scale any part of our API, giving us a huge amount of flexibility in performance tuning. AppSync allows us to define regular invoke and batch invoke resolvers, allowing us to write lambda functions that act as dataloaders, avoiding the N+1 query problem.

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

A lambda function serving the `createItem` mutation. Returns [`responses.CreateItem`](internal/responses/create_item.go) using the `Name` from [`requests.CreateItem`](internal/requests/create_item.go).

| meta    | value                                                                |
| ------- | -------------------------------------------------------------------- |
| runtime | go1.x                                                                |
| uri     | [cmd/create-item/main.go](cmd/create-item/main.go)                   |
| handler | [internal/handlers/create_item.go](internal/handlers/create_item.go) |

### ReadItemFunction

A lambda function serving the `item` query. Returns [`responses.Item`](internal/responses/item.go).

#### ⛔️ Note: Currently returns a not implemented error.

| meta    | value                                                            |
| ------- | ---------------------------------------------------------------- |
| runtime | go1.x                                                            |
| uri     | [cmd/read-item/main.go](cmd/read-item/main.go)                   |
| handler | [internal/handlers/read_item.go](internal/handlers/read_item.go) |

### ListItemsFunction

A lambda function serving the `items` query. Returns [`responses.Item[]`](internal/responses/item.go). This handler is invoked as a batch invoke and operates in a similar fashion to a data loader.

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
