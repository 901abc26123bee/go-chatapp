db.createUser(
  {
    user: "mongouser",
    pwd: "mongopw",
    roles: [
      {
        role: "readWrite",
        db: "gsm"
      }
    ]
  }
);