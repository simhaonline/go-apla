// Copyright (C) 2017, 2018, 2019 EGAAS S.A.
//
// This program is free software; you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation; either version 2 of the License, or (at
// your option) any later version.
//
// This program is distributed in the hope that it will be useful, but
// WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
// General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program; if not, write to the Free Software
// Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA 02110-1301, USA.

package migration

var firstSystemParametersDataSQL = `
INSERT INTO "1_system_parameters" ("id","name", "value", "conditions") VALUES 
	('1','default_ecosystem_page', 'If(#ecosystem_id# > 1){Include(@1welcome)}', 'ContractAccess("@1UpdateSysParam")'),
	('2','default_ecosystem_menu', '', 'ContractAccess("@1UpdateSysParam")'),
	('3','default_ecosystem_contract', '', 'ContractAccess("@1UpdateSysParam")'),
	('4','gap_between_blocks', '2', 'ContractAccess("@1UpdateSysParam")'),
	('5','rollback_blocks', '60', 'ContractAccess("@1UpdateSysParam")'),
	('6','new_version_url', 'upd.apla.io', 'ContractAccess("@1UpdateSysParam")'),
	('7','full_nodes', '', 'ContractAccess("@1UpdateSysParam","@1NodeRemoveByKey")'),
	('8','number_of_nodes', '101', 'ContractAccess("@1UpdateSysParam")'),
	('9','price_create_contract', '200', 'ContractAccess("@1UpdateSysParam")'),
	('10','price_create_menu', '100', 'ContractAccess("@1UpdateSysParam")'),
	('11','price_create_page', '100', 'ContractAccess("@1UpdateSysParam")'),
	('12','blockchain_url', '', 'ContractAccess("@1UpdateSysParam")'),
	('13','max_block_size', '67108864', 'ContractAccess("@1UpdateSysParam")'),
	('14','max_tx_size', '33554432', 'ContractAccess("@1UpdateSysParam")'),
	('15','max_tx_block', '1000', 'ContractAccess("@1UpdateSysParam")'),
	('16','max_columns', '50', 'ContractAccess("@1UpdateSysParam")'),
	('17','max_indexes', '5', 'ContractAccess("@1UpdateSysParam")'),
	('18','max_tx_block_per_user', '100', 'ContractAccess("@1UpdateSysParam")'),
	('19','max_fuel_tx', '20000000', 'ContractAccess("@1UpdateSysParam")'),
	('20','max_fuel_block', '200000000', 'ContractAccess("@1UpdateSysParam")'),
	('21','commission_size', '3', 'ContractAccess("@1UpdateSysParam")'),
	('22','commission_wallet', '', 'ContractAccess("@1UpdateSysParam")'),
	('23','fuel_rate', '[["1","100000000000"]]', 'ContractAccess("@1UpdateSysParam")'),
	('24','price_exec_address_to_id', '10', 'ContractAccess("@1UpdateSysParam")'),
	('25','price_exec_id_to_address', '10', 'ContractAccess("@1UpdateSysParam")'),
	('26','price_exec_sha256', '50', 'ContractAccess("@1UpdateSysParam")'),
	('27','price_exec_pub_to_id', '10', 'ContractAccess("@1UpdateSysParam")'),
	('28','price_exec_ecosys_param', '10', 'ContractAccess("@1UpdateSysParam")'),
	('29','price_exec_sys_param_string', '10', 'ContractAccess("@1UpdateSysParam")'),
	('30','price_exec_sys_param_int', '10', 'ContractAccess("@1UpdateSysParam")'),
	('31','price_exec_sys_fuel', '10', 'ContractAccess("@1UpdateSysParam")'),
	('32','price_exec_validate_condition', '30', 'ContractAccess("@1UpdateSysParam")'),
	('33','price_exec_eval_condition', '20', 'ContractAccess("@1UpdateSysParam")'),
	('34','price_exec_has_prefix', '10', 'ContractAccess("@1UpdateSysParam")'),
	('35','price_exec_contains', '10', 'ContractAccess("@1UpdateSysParam")'),
	('36','price_exec_replace', '10', 'ContractAccess("@1UpdateSysParam")'),
	('37','price_exec_join', '10', 'ContractAccess("@1UpdateSysParam")'),
	('38','price_exec_update_lang', '10', 'ContractAccess("@1UpdateSysParam")'),
	('39','price_exec_size', '10', 'ContractAccess("@1UpdateSysParam")'),
	('40','price_exec_substr', '10', 'ContractAccess("@1UpdateSysParam")'),
	('41','price_exec_contracts_list', '10', 'ContractAccess("@1UpdateSysParam")'),
	('42','price_exec_is_object', '10', 'ContractAccess("@1UpdateSysParam")'),
	('43','price_exec_compile_contract', '100', 'ContractAccess("@1UpdateSysParam")'),
	('44','price_exec_flush_contract', '50', 'ContractAccess("@1UpdateSysParam")'),
	('45','price_exec_eval', '10', 'ContractAccess("@1UpdateSysParam")'),
	('46','price_exec_len', '5', 'ContractAccess("@1UpdateSysParam")'),
	('47','price_exec_bind_wallet', '10', 'ContractAccess("@1UpdateSysParam")'),
	('48','price_exec_unbind_wallet', '10', 'ContractAccess("@1UpdateSysParam")'),
	('49','price_exec_create_ecosystem', '100', 'ContractAccess("@1UpdateSysParam")'),
	('50','price_exec_table_conditions', '100', 'ContractAccess("@1UpdateSysParam")'),
	('51','price_exec_create_table', '100', 'ContractAccess("@1UpdateSysParam")'),
	('52','price_exec_perm_table', '100', 'ContractAccess("@1UpdateSysParam")'),
	('53','price_exec_column_condition', '50', 'ContractAccess("@1UpdateSysParam")'),
	('54','price_exec_create_column', '50', 'ContractAccess("@1UpdateSysParam")'),
	('55','price_exec_perm_column', '50', 'ContractAccess("@1UpdateSysParam")'),
	('56','price_exec_json_to_map', '50', 'ContractAccess("@1UpdateSysParam")'),
	('57','max_block_generation_time', '2000', 'ContractAccess("@1UpdateSysParam")'),
	('58','block_reward','1000','ContractAccess("@1UpdateSysParam")'),
	('59','incorrect_blocks_per_day','10','ContractAccess("@1UpdateSysParam")'),
	('60','node_ban_time','86400000','ContractAccess("@1UpdateSysParam")'),
	('61','local_node_ban_time','1800000','ContractAccess("@1UpdateSysParam")'),
	('62','test','false','false'),
	('63','price_tx_data', '0', 'ContractAccess("@1UpdateSysParam")'),
	('64', 'price_exec_contract_by_name', '0', 'ContractAccess("@1UpdateSysParam")'),
	('65', 'price_exec_contract_by_id', '0', 'ContractAccess("@1UpdateSysParam")'),
	('66','private_blockchain', '1', 'false');
`
