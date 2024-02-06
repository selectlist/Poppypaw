import * as database from "../../Serendipity/prisma.js";
import { FastifyReply, FastifyRequest } from "fastify";

export default {
	method: "POST",
	url: "/users/applications",
	schema: {
		summary: "Create an Developer Application",
		description:
			"Returns boolean value indicating whether the creation was successful or not.",
		tags: ["developers"],
		body: {
			type: "object",
			properties: {
				name: { type: "string" },
				logo: { type: "string" },
			},
			required: ["name", "logo"],
		},
		security: [
			{
				apiKey: [],
			},
		],
	},
	handler: async (request: FastifyRequest, reply: FastifyReply) => {
		try {
			const { name, logo }: any = request.body;
			const Authorization: any = request.headers.authorization;
			const token = await database.Tokens.get({ token: Authorization });
			const dbUser = await database.Users.get({ userid: token.userid });

			if (dbUser) {
				const apps = await database.Applications.createApp(
					dbUser?.userid,
					name,
					logo
				);

				return reply.send(apps);
			} else
				return reply.send({
					token: Authorization,
					error: true,
					message: "User does not exist.",
				});
		} catch (error) {
			reply.status(500).send({
				error: "Internal Server Error",
				message: error.errorInfo.message,
			});
		}
	},
};
