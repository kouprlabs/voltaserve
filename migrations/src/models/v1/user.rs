use sea_orm_migration::prelude::*;

#[derive(Iden)]
pub enum User {
    Table,
    Id,
    FullName,
    Username,
    Email,
    PasswordHash,
    RefreshTokenValue,
    RefreshTokenExpiry,
    ResetPasswordToken,
    EmailConfirmationToken,
    EmailUpdateToken,
    EmailUpdateValue,
    IsEmailConfirmed,
    Picture,
    CreateTime,
    UpdateTime,
}