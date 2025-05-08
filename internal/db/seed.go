package db

import (
	"context"
	"fmt"
	"log"
	"math/rand"

	"github.com/puremike/social-go/internal/model"
	"github.com/puremike/social-go/internal/store"
)

var usernames = []string{"ShadowWolf99", "NeonPhoenix23", "CyberPunk87", "MysticDragon45", "QuantumRider12", "SilverFox34", "CrimsonTiger78", "ElectricEagle56", "GoldenLion89", "NightOwl21", "SteelHawk43", "FrostGiant67", "IronWolf98", "SolarFlare54", "ThunderBear76", "OceanMaster32", "DarkKnight88", "BlazeRunner19", "StormChaser47", "LunarWolf65", "CosmicDrifter29", "InfernoDragon73", "TitaniumFox81", "PhantomRider36", "ArcticWolf58", "GalacticEagle92", "CrimsonPhoenix14", "ShadowHunter77", "NeonTiger63", "CyberWolf48", "MysticEagle25", "QuantumFox91", "SilverDragon39", "GoldenHawk72", "NightWolf84", "SteelTiger16", "FrostWolf53", "IronEagle28", "SolarTiger71", "ThunderWolf95", "OceanEagle37", "DarkFox82", "BlazeWolf64", "StormWolf49", "LunarEagle18", "CosmicWolf57", "InfernoFox86", "TitaniumWolf33", "PhantomFox79", "ArcticEagle44", "GalacticWolf69", "CrimsonWolf27", "ShadowFox59", "NeonEagle74", "CyberDragon13", "MysticWolf68", "QuantumEagle22", "SilverTiger88", "GoldenWolf41", "NightEagle55", "SteelDragon93", "FrostEagle17", "IronTiger84", "SolarWolf46", "ThunderEagle72", "OceanWolf38", "DarkEagle91", "BlazeEagle63", "StormEagle29", "LunarTiger75", "CosmicEagle51", "InfernoWolf66", "TitaniumEagle34", "PhantomWolf87", "ArcticTiger58", "GalacticEagle92", "CrimsonFox14", "ShadowEagle77", "NeonWolf63", "CyberTiger48", "MysticEagle25", "QuantumFox91", "SilverDragon39", "GoldenHawk72", "NightWolf84", "SteelTiger16", "FrostWolf53", "IronEagle28", "SolarTiger71", "ThunderWolf95", "OceanEagle37", "DarkFox82", "BlazeWolf64", "StormWolf49", "LunarEagle18", "CosmicWolf57", "InfernoFox86", "TitaniumWolf33", "PhantomFox79", "ArcticEagle44", "GalacticWolf69"}

var titles = []string{
	"10 Tips for Boosting Your Productivity Today", "The Future of Artificial Intelligence: What to Expect", "How to Travel on a Budget Without Sacrificing Comfort", "Exploring the Hidden Gems of Europe", "The Science Behind Healthy Eating Habits", "Top 5 Gadgets You Need in 2024", "Why Mindfulness is the Key to a Balanced Life", "The Rise of Remote Work: Pros and Cons", "How to Build a Successful Side Hustle from Scratch", "The Best Books to Read for Personal Growth", "Understanding Blockchain Technology in Simple Terms", "How to Stay Motivated When Pursuing Long-Term Goals", "The Impact of Social Media on Mental Health", "A Beginner's Guide to Investing in the Stock Market", "The Art of Minimalism: Simplifying Your Life", "How to Master the Skill of Public Speaking", "The Role of Renewable Energy in Combating Climate Change", "Exploring the Wonders of Space: A Beginner's Guide", "How to Build Stronger Relationships in a Digital Age", "The Benefits of Learning a New Language in 2024"}

var contents = []string{
	"Productivity doesn't have to be complicated. Start by organizing your tasks, setting clear goals, and eliminating distractions. Small changes like these can make a huge difference in your daily output.", "Artificial intelligence is evolving rapidly, and its applications are endless. From healthcare to finance, AI is transforming industries and creating new opportunities for innovation.", "Traveling on a budget doesn't mean you have to compromise on comfort. With careful planning, smart booking, and a bit of creativity, you can enjoy amazing experiences without breaking the bank.", "Europe is full of hidden gems waiting to be discovered. From quaint villages to stunning landscapes, there's so much more to explore beyond the usual tourist hotspots.", "Healthy eating is more than just a trend—it's a lifestyle. Understanding the science behind nutrition can help you make better choices and improve your overall well-being.", "Technology is advancing faster than ever, and these five gadgets are must-haves for 2024. Stay ahead of the curve with tools that make life easier and more efficient.", "Mindfulness is more than just meditation. It's about being present in the moment and cultivating a sense of calm and clarity in your everyday life.", "Remote work has become the new norm, but it comes with its own set of challenges. Learn how to navigate the pros and cons to make the most of this flexible work style.", "Starting a side hustle can be daunting, but with the right mindset and strategy, you can turn your passion into a profitable venture. Here's how to get started.", "Books have the power to transform your life. Whether you're looking for inspiration or practical advice, these titles are a great place to start your personal growth journey.", "Blockchain technology is revolutionizing the way we think about data and transactions. This guide breaks down the basics so you can understand its potential.", "Staying motivated over the long haul can be tough. Setting small milestones, celebrating progress, and staying focused on your 'why' can keep you on track.", "Social media has a profound impact on mental health. While it connects us, it can also lead to comparison and anxiety. Learn how to use it mindfully.", "Investing in the stock market can seem intimidating, but it doesn't have to be. Start with the basics, do your research, and take it one step at a time.", "Minimalism is about more than just decluttering. It's a mindset that helps you focus on what truly matters and let go of the rest.", "Public speaking is a skill that can be learned and mastered. With practice and confidence, you can captivate any audience and deliver your message effectively.", "Renewable energy is key to a sustainable future. From solar to wind power, these technologies are paving the way for a cleaner, greener planet.", "Space exploration continues to fascinate us. Whether you're a beginner or an enthusiast, there's always something new to learn about the universe.", "In a world dominated by screens, building meaningful relationships can be challenging. Here are some tips to foster deeper connections in the digital age.", "Learning a new language opens doors to new cultures and opportunities. In 2024, make it a goal to expand your horizons and embrace the benefits of bilingualism."}

var tags = []string{
	"productivity", "AI", "travel", "Europe", "health", "gadgets", "mindfulness", "remote work", "side hustle", "books", "blockchain", "motivation", "social media", "investing", "minimalism", "public speaking", "renewable energy", "space", "relationships", "language learning"}

var comments = []string{
	"This is such an insightful post! I learned so much.", "I never thought about it this way. Thanks for sharing!", "Great tips! I’ll definitely try these out.", "This is exactly what I needed to read today.", "I have a different perspective, but this was still interesting.",
	"Can you recommend any resources to learn more about this?", "This post inspired me to take action. Thank you!", "I’ve been struggling with this, and your advice really helped.", "Such a well-written and thought-provoking article.", "I love how practical and actionable your advice is.", "This is so relatable! Thanks for putting it into words.", "I’ve been looking for something like this. Great job!", "Your post made me rethink my approach. Much appreciated!", "I’m sharing this with my friends—it’s too good not to!", "This is a game-changer. Thank you for the insights!", "I’ve tried this before, and it really works. Highly recommend!", "Your writing style is so engaging. Keep it up!", "This is a fresh take on a topic I thought I knew well.", "I’m excited to implement these ideas. Thanks for the inspiration!", "This post is a goldmine of information. Well done!"}

func Seed(store *store.Storage) {
	ctx := context.Background()

	users := generateUsers(100)
	for _, user := range users {
		if err := store.Users.Create(ctx, user); err != nil {
			log.Println("Error creating user: ", err)
			return
		}
	}

	posts := generatePosts(100, users)
	for _, post := range posts {
		if err := store.Posts.Create(ctx, post); err != nil {
			log.Println("Error creating posts: ", err)
			return
		}
	}

	comments := generateComments(100, users, posts)
	for _, comment := range comments {
		if err := store.Comments.Create(ctx, comment); err != nil {
			log.Println("Error creating comment:", err)
			return
		}
	}

	log.Println("Seeding complete")
}

func generateUsers(num int) []*store.UserModel {

	users := make([]*store.UserModel, num)
	for i := 0; i < num; i++ {
		users[i] = &store.UserModel{
			Username: usernames[i%len(usernames)] + fmt.Sprintf("%d", i),
			Email:    usernames[i%len(usernames)] + fmt.Sprintf("%d", i) + "@email.com",
			RoleId:   1,
		}
	}
	return users

}

func generatePosts(num int, users []*store.UserModel) []*model.PostModel {
	posts := make([]*model.PostModel, num)

	for i := 0; i < num; i++ {
		user := users[rand.Intn((len(users)))]

		posts[i] = &model.PostModel{
			UserID:  user.ID,
			Title:   titles[rand.Intn(len(titles))],
			Content: contents[rand.Intn(len(contents))],
			Tags: []string{
				tags[rand.Intn(len(tags))], tags[rand.Intn(len(tags))],
			},
		}
	}
	return posts
}

func generateComments(num int, users []*store.UserModel, posts []*model.PostModel) []*model.CommentModel {
	cms := make([]*model.CommentModel, num)
	for i := 0; i < num; i++ {
		cms[i] = &model.CommentModel{
			PostID:  posts[rand.Intn(len(posts))].ID,
			UserID:  users[rand.Intn(len(users))].ID,
			Content: comments[rand.Intn(len(comments))],
		}
	}
	return cms
}
