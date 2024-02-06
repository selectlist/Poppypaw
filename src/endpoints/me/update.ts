import { FastifyReply, FastifyRequest } from "fastify";
import * as database from "../../Serendipity/prisma.js";

export default {
	url: "/users/@me",
	method: "PATCH",
	schema: {
		summary: "Update @me information",
		description:
			"Returns boolean value indicating whether the update was successful or not.",
		tags: ["users"],
		body: {
			type: "object",
			properties: {
				bio: { type: "string" },
			},
			required: ["bio"],
		},
		security: [
			{
				apiKey: [],
			},
		],
	},
	handler: async (request: FastifyRequest, reply: FastifyReply) => {
		let data = request.body;
		const Authorization: any = request.headers.authorization;
		const token = await database.Tokens.get({ token: Authorization });
		const user = await database.Users.get({ userid: token.userid });

		if (user) {
			if (!data["bio"] || data["bio"] === "") data["bio"] = null;

			await database.Users.update(user.userid, {
				username: user.username,
				userid: user.userid,
				revoltid: user.revoltid,
				bio: data["bio"],
				avatar: user.avatar,
				badges: user.badges,
				staff_perms: user.staff_perms,
			});

			return reply.send({
				success: true,
			});
		} else
			reply.status(404).send({
				message:
					"We couldn't fetch any information about you in our database",
				token: Authorization,
				error: true,
			});
	},
};
