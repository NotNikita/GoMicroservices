db = db.getSiblingDB("udemy_logs");

// Create collections
db.createCollection("logs");

// Insert seed data if collection is empty
if (db.logs.countDocuments() === 0) {
  db.logs.insertMany([
    {
      name: "Event",
      data: {
        event: "database_seeded",
        service: "logger",
        timestamp: new Date(),
      },
    },
  ]);
  print("Database seeded successfully");
} else {
  print("Database already contains data, skipping seed");
}
