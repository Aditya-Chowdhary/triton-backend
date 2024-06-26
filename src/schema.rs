table! {
    pastes (id) {
        id -> Nullable<Varchar>,
        belongs_to -> Nullable<Varchar>,
        is_url -> Nullable<Bool>,
        content -> Text,
    }
}

table! {
    users (id) {
        id -> Varchar,
        username -> Nullable<Varchar>,
        password -> Nullable<Varchar>,
        activated -> Nullable<Bool>,
    }
}

joinable!(pastes -> users (belongs_to));

allow_tables_to_appear_in_same_query!(pastes, users,);
