use anyhow::Result;
use diesel::pg::PgConnection;
use diesel::prelude::*;

use crate::core::users::entity::User;
use crate::schema::users;

pub fn create_user(user: &User, conn: &PgConnection) -> Result<usize> {
    let records = diesel::insert_into(users::table)
        .values(user)
        .on_conflict_do_nothing()
        .execute(conn)?;
    Ok(records)
}

pub fn find_user(id: String, conn: &PgConnection) -> Result<User> {
    let user = users::table.find(id).get_result::<User>(conn)?;
    Ok(user)
}

pub fn update_user(user: &User, conn: &PgConnection) -> Result<User> {
    let user = diesel::update(users::table.find(user.id.clone()))
        .set(user)
        .get_result::<User>(conn)?;
    Ok(user)
}
