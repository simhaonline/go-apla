package migration

var tablesDataSQL = `INSERT INTO "%[1]d_tables" ("id", "name", "permissions","columns", "conditions") VALUES 
('1', 'contracts', '{"insert": "ContractConditions(\"MainCondition\")", "update": "ContractConditions(\"MainCondition\")", "new_column": "ContractConditions(\"MainCondition\")"}', 
'{"name": "false", 
	"value": "ContractConditions(\"MainCondition\")",
	  "wallet_id": "ContractConditions(\"MainCondition\")",
	  "token_id": "ContractConditions(\"MainCondition\")",
		"active": "ContractConditions(\"MainCondition\")",
		"confirmation": "ContractConditions(\"MainCondition\")",
	  "conditions": "ContractConditions(\"MainCondition\")"}', 'ContractAccess("@1EditTable")'),
	('2', 'keys', 
	'{"insert": "true", "update": "true", 
	  "new_column": "ContractConditions(\"MainCondition\")"}',
	'{"pub": "ContractConditions(\"MainCondition\")",
	  "amount": "ContractConditions(\"MainCondition\")",
	  "maxpay": "ContractConditions(\"MainCondition\")",
	  "deleted": "ContractConditions(\"MainCondition\")",
	  "blocked": "ContractConditions(\"MainCondition\")",
	  "multi": "ContractAccess(\"MultiwalletCreate\")"}', 
	'ContractAccess("@1EditTable")'),
	('3', 'history', 
	'{"insert": "ContractConditions(\"NodeOwnerCondition\")", "update": "ContractConditions(\"MainCondition\")", 
	  "new_column": "ContractConditions(\"MainCondition\")"}',
	'{"sender_id": "ContractConditions(\"MainCondition\")",
	  "recipient_id": "ContractConditions(\"MainCondition\")",
	  "amount":  "ContractConditions(\"MainCondition\")",
	  "comment": "ContractConditions(\"MainCondition\")",
	  "block_id":  "ContractConditions(\"MainCondition\")",
	  "txhash": "ContractConditions(\"MainCondition\")"}', 'ContractAccess("@1EditTable")'),        
	('4', 'languages', 
	'{"insert": "ContractConditions(\"MainCondition\")", "update": "ContractConditions(\"MainCondition\")", 
	  "new_column": "ContractConditions(\"MainCondition\")"}',
	'{"app_id": "ContractConditions(\"MainCondition\")",
	  "name": "ContractConditions(\"MainCondition\")",
	  "res": "ContractConditions(\"MainCondition\")",
	  "conditions": "ContractConditions(\"MainCondition\")",
	  "app_id": "ContractConditions(\"MainCondition\")"}', 'ContractAccess("@1EditTable")'),
	('5', 'menu', 
		'{"insert": "ContractConditions(\"MainCondition\")", "update": "ContractConditions(\"MainCondition\")", 
	  "new_column": "ContractConditions(\"MainCondition\")"}',
	'{"name": "ContractConditions(\"MainCondition\")",
"value": "ContractConditions(\"MainCondition\")",
"conditions": "ContractConditions(\"MainCondition\")"
	}', 'ContractAccess("@1EditTable")'),
	('6', 'pages', 
		'{"insert": "ContractConditions(\"MainCondition\")", "update": "ContractConditions(\"MainCondition\")", 
	  "new_column": "ContractConditions(\"MainCondition\")"}',
	'{"name": "ContractConditions(\"MainCondition\")",
"value": "ContractConditions(\"MainCondition\")",
"menu": "ContractConditions(\"MainCondition\")",
"validate_count": "ContractConditions(\"MainCondition\")",
"validate_mode": "ContractConditions(\"MainCondition\")",
"app_id": "ContractConditions(\"MainCondition\")",
"conditions": "ContractConditions(\"MainCondition\")"
	}', 'ContractAccess("@1EditTable")'),
	('7', 'blocks', 
	'{"insert": "ContractConditions(\"MainCondition\")", "update": "ContractConditions(\"MainCondition\")", 
	  "new_column": "ContractConditions(\"MainCondition\")"}',
	'{"name": "ContractConditions(\"MainCondition\")",
"value": "ContractConditions(\"MainCondition\")",
"conditions": "ContractConditions(\"MainCondition\")"
	}', 'ContractAccess("@1EditTable")'),
	('9', 'members', 
		'{"insert":"ContractAccess(\"Profile_Edit\")","update":"true","new_column":"ContractConditions(\"MainCondition\")"}',
		'{"image_id":"ContractAccess(\"ProfileEditAvatar\")","member_info":"ContractAccess(\"Profile_Edit\")","member_name":"false"}', 
		'ContractConditions("MainCondition")'),
	('10', 'roles',
		'{"insert":"ContractAccess(\"Roles_Create\")",
			"update":"ContractConditions(\"MainCondition\")",
			"new_column":"ContractConditions(\"MainCondition\")"}', 
		'{"default_page":"false",
			"creator":"false",
			"deleted":"ContractAccess(\"Roles_Del\")",
			"company_id":"false",
			"date_deleted":"ContractAccess(\"Roles_Del\")",
			"image_id":"ContractAccess(\"Roles_Create\")",
			"role_name":"false",
			"date_created":"false",
			"roles_access":"ContractAccess(\"Roles_AccessManager\")",
			"role_type":"false"}',
		'ContractConditions("MainCondition")'),
	('11', 'roles_participants',
		'{"insert":"ContractAccess(\"Roles_Assign\",\"voting_CheckDecision\")",
			"update":"ContractConditions(\"MainCondition\")",
			"new_column":"ContractConditions(\"MainCondition\")"}',
		'{"deleted":"ContractAccess(\"Roles_Unassign\")",
			"date_deleted":"ContractAccess(\"Roles_Unassign\")",
			"member":"false",
			"role":"false",
			"date_created":"false",
			"appointed":"false"}', 
		'ContractConditions("MainCondition")'),
	('12', 'notifications',
		'{"insert":"ContractAccess(\"notifications_Send\", \"CheckNodesBan\")",
			"update":"ContractAccess(\"notifications_Send\", \"notifications_Close\", \"notifications_Process\")",
			"new_column":"ContractConditions(\"MainCondition\")"}',
		'{"date_closed":"ContractAccess(\"notifications_Close\")",
			"sender":"false",
			"processing_info":"ContractAccess(\"notifications_Close\",\"notifications_Process\")",
			"date_start_processing":"ContractAccess(\"notifications_Close\",\"notifications_Process\")",
			"notification":"false",
			"page_name":"false",
			"page_params":"false",
			"closed":"ContractAccess(\"notifications_Close\")",
			"date_created":"false",
			"recipient":"false"}',
		'ContractAccess("@1EditTable")'),
	('13', 'sections', 
		'{"insert": "ContractConditions(\"MainCondition\")", "update": "ContractConditions(\"MainCondition\")", 
		"new_column": "ContractConditions(\"MainCondition\")"}',
		'{"title": "ContractConditions(\"MainCondition\")",
			"urlname": "ContractConditions(\"MainCondition\")",
			"page": "ContractConditions(\"MainCondition\")",
			"roles_access": "ContractConditions(\"MainCondition\")",
			"status": "ContractConditions(\"MainCondition\")"}', 
			'ContractConditions("MainCondition")'),
	('14', 'applications',
		'{"insert": "ContractConditions(\"MainCondition\")",
			 "update": "ContractConditions(\"MainCondition\")", 
			 "new_column": "ContractConditions(\"MainCondition\")"}',
		'{"name": "ContractConditions(\"MainCondition\")",
		  "uuid": "false",
		  "conditions": "ContractConditions(\"MainCondition\")",
		  "deleted": "ContractConditions(\"MainCondition\")"}',
		'ContractConditions("MainCondition")'),
	('15', 'binaries',
		'{"insert":"ContractAccess(\"@1UploadBinary\")",
			"update":"ContractAccess(\"@1UploadBinary\")",
			"new_column":"ContractConditions(\"MainCondition\")"}',
		'{"hash":"ContractAccess(\"@1UploadBinary\")",
			"member_id":"false",
			"data":"ContractAccess(\"@1UploadBinary\")",
			"name":"false",
			"app_id":"false"}',
		'ContractConditions(\"MainCondition\")'),
	('16', 'parameters',
		'{"insert": "ContractConditions(\"MainCondition\")", "update": "ContractConditions(\"MainCondition\")",
			"new_column": "ContractConditions(\"MainCondition\")"}',
		'{"name": "ContractConditions(\"MainCondition\")",
			"value": "ContractConditions(\"MainCondition\")",
			"conditions": "ContractConditions(\"MainCondition\")"}',
		'ContractAccess("@1EditTable")'),
	('17', 'app_params',
		'{"insert": "ContractConditions(\"MainCondition\")", "update": "ContractConditions(\"MainCondition\")",
			"new_column": "ContractConditions(\"MainCondition\")"}',
		'{"app_id": "ContractConditions(\"MainCondition\")",
			"name": "ContractConditions(\"MainCondition\")",
			"value": "ContractConditions(\"MainCondition\")",
			"conditions": "ContractConditions(\"MainCondition\")"}',
		'ContractAccess("@1EditTable")'),
	('19', 'buffer_data',
		'{"insert":"true","update":"true",
			"new_column":"ContractConditions(\"MainCondition\")"}',
		'{"key": "false",
			"value": "true",
			"member_id": "false"}',
		'ContractConditions("MainCondition")');
`
