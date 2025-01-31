# appsync-lambda-golang

A sample AWS AppSync GraphQL API leveraging lambda functions as resolvers.

## Why?

GraphQL (when implemented correctly) is a fantastic way to design highly expressive API's. 

With a bit of additional configuration, AWS AppSync allows us to define Lambda functions as resolvers, unlocking a new level of customization and scalability.

Using this approach we can independently develop, deploy, and scale any resolver in our API, giving us a huge amount of flexibility in how we design and maintain our system.

And AWS AppSync allows us to perform regular and batch invocation of lambda resolvers, allowing us to write lambda functions that avoid the N + 1 query problem.

## GraphQL Schema

This sample GraphQL API has the following schema.

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

## Lambda Functions

### CreateItemFunction

Resolver for `createItem` mutation. Returns [`responses.CreateItem`](internal/responses/create_item.go) using the `Name` from [`requests.CreateItem`](internal/requests/create_item.go).

| meta    | value                                                                |
| ------- | -------------------------------------------------------------------- |
| runtime | go1.x                                                                |
| uri     | [cmd/create-item/main.go](cmd/create-item/main.go)                   |
| handler | [internal/resolvers/create_item.go](internal/resolvers/create_item.go) |

### ReadItemFunction

Resolver for `item` query. Returns [`responses.Item`](internal/responses/item.go).

#### ⛔️ Note: Currently returns a not implemented error.

| meta    | value                                                            |
| ------- | ---------------------------------------------------------------- |
| runtime | go1.x                                                            |
| uri     | [cmd/read-item/main.go](cmd/read-item/main.go)                   |
| handler | [internal/resolvers/read_item.go](internal/resolvers/read_item.go) |

### ListItemsFunction

Resolver for `items` query. Returns [`responses.Item[]`](internal/responses/item.go). This handler is invoked as a batch invoke and operates in a similar fashion to a data loader.

#### ⛔️ Note: Currently returns a not implemented error.

| meta    | value                                                              |
| ------- | ------------------------------------------------------------------ |
| runtime | go1.x                                                              |
| uri     | [cmd/list-items/main.go](cmd/list-items/main.go)                   |
| handler | [internal/resolvers/list_items.go](internal/resolvers/list_items.go) |

## Building

`sam build`

## Deploying

Initial: `sam deploy --guided`

Subsequent: `sam deploy`
