// Copyright (c) 2018-Now Dunyu All Rights Reserved.
//
// Author      : https://www.wengold.net
// Email       : support@wengold.net
//
// Prismy.No | Date       | Modified by. | Description
// -------------------------------------------------------------------
// 00001       2019/05/22   yangping       New version
// -------------------------------------------------------------------

package invar

// Fixed target payee for all trade of DY
const PayeeDunYu = "PAYEE_DUNYU"

// Trade status machine:
// ===================================================================
//
//                      + --------------------------------------- +
//    + ---------------/------- +                                 |
//    |               /         |                      [ END ]    |
//    |              /          v                            \    v
//    |     TSUnpaid --- + -> TSRevoked --------------------> TSClosed
//    |     [ TRADE ]    |      ^                                 ^
//    |                  |      |                                 |
//    |                  + -> TSPayError <- + (max counts 5) ---- +
//    |                  |      |           |                     |
//    |                  |      + --------- +                     |
//    |                  |      v                                 |
//    |                  + -> TSPaid -------- + ----------------- +
//    |                  |    [ SUCCESS ]     ^                   |
//    |                  |                    |                   |
//    |                  + -> TSCompleted --- +                   |
//    |                                                           |
//    |     [ REFUND ]      [ SUCCESS ]                           |
//    + --- TSInProgress -> TSRefund ---------------------------- +
//              |              ^                                  |
//              |              |                                  |
//              + ------- > TSRefundError <- + (max counts 5) --- +
//              |              |             |                    |
//              |              + ----------- +                    |
//              + ----------------------------------------------- +
//
// [ TRADE   ] : Trade  transaction start and default state
// [ REFUND  ] : Refund transaction start and default state
// [ SUCCESS ] : Success paid or refund status
// [ END     ] : Closed ticket state
// ===================================================================

// Unpid state, can be use as default trade state for generate
// a trade ticket, as status machine it can change to :
//
//	TSUnpaid	 -> TSRevoked   : canceld parment
//				 -> TSPayError  : pay error
//				 -> TSPaid      : success paid
//				 -> TSCompleted : only for dividing payment, to mark dividing completed
const TSUnpaid = "UNPAID"

// Pay error state, as status machine it can change to :
//
//	TSPayError	 -> TSPaid    : success paid
//				 -> TSRevoked : canceld parment
//				 -> self (over 5 times) -> TSClosed
const TSPayError = "PAY_ERROR"

// Revoked state, cancel by user, as status machine it only
// can be changed to closed state:
//
//	TSRevoked	 -> TSClosed : close the trade ticket
const TSRevoked = "REVOKED"

// Paied success state, as status machine it only can be
// changed to closed state :
//
//	TSPaid		 -> TSClosed : close the trade ticket
//
// `WARNING` :
//
// the refund action will generate a new trade ticket and set
// TSInProgress as default.
const TSPaid = "PAID"

// Completed all dividing payments, as status machine it only
// can be changed to closed state :
//
//	TSCompleted	 -> TSClosed : close the trade ticket
const TSCompleted = "COMPLETED"

// Refund in progress state, use as default trade state when generate
// a refund ticket, as status machine it can change to :
//
//	TSInProgress -> TSRefund      : success refund
//				 -> TSRefundError : refund error
const TSInProgress = "REFUND_IN_PROGRESS"

// Refund success state, as status machine it only can be
// changed to closed state :
//
//	TSRefundError	-> TSClosed : close the trade ticket
//					-> self (over 5 times) -> TSClosed
const TSRefundError = "REFUND_ERROR"

// Refund success state, as status machine it only can be
// changed to closed state :
//
//	TSRefund	 -> TSClosed : close the trade ticket
const TSRefund = "REFUND"

// Closed state, as the last state of status machine.
const TSClosed = "CLOSED"
