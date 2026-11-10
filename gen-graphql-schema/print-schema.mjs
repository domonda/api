// Prints the schema of a GraphQL endpoint in SDL to stdout,
// built from an introspection query, which the domonda API
// answers without an API key.
import { buildClientSchema, getIntrospectionQuery, printSchema } from "graphql";

const endpoint = process.argv[2] ?? "https://domonda.app/api/public/graphql";

const response = await fetch(endpoint, {
  method: "POST",
  headers: { "Content-Type": "application/json" },
  body: JSON.stringify({ query: getIntrospectionQuery() }),
});
if (!response.ok) {
  throw new Error(`${endpoint} responded ${response.status} ${response.statusText}`);
}
const { data, errors } = await response.json();
if (errors?.length) {
  throw new Error(`introspection failed: ${JSON.stringify(errors)}`);
}
process.stdout.write(printSchema(buildClientSchema(data)) + "\n");
