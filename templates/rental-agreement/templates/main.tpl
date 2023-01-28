{{define "simple-paragraph"}}
{{$head_text := .Head}}
{{$text := .Content}}
<paragraph margin="20 0 0 0">
<text-chunk font="times-bold" font-size="13">{{$head_text}}</text-chunk>
<text-chunk font="times" font-size="13">{{$text}}</text-chunk>
</paragraph>
{{end}}

{{ define "form-sig"}}
   {{$margin := .Margin}}
   {{$text := .Text}}
   <division>
      <paragraph>
         <text-chunk font= "times-bold" font-size="11">{{$text}}: </text-chunk>
      </paragraph>
      <line fit-mode="fill-width" position="relative" thickness= "0.5" margin="{{$margin}}"></line>
   </division>
{{end}}
{{ define "simple-form"}}
   {{$margin := .Margin}}
   {{$text := .Text}}
   <division>
      <paragraph>
         <text-chunk font= "times" font-size="11">{{$text}} </text-chunk>
      </paragraph>
      <line fit-mode="fill-width" position="relative" thickness= "0.5" margin="{{$margin}}"></line>
   </division>
{{end}}
<paragraph margin="0 0 0 0" text-align = "center">
   <text-chunk font="times-bold" font-size="21.5"> LEASE WITH OPTION TO PURCHASE </text-chunk>
</paragraph>

<paragraph margin="20 0 0 0">
<text-chunk font="times" font-size="13">This agreement, dated December 9 2020, by and between a business entity known as </text-chunk>
<text-chunk font="times" font-size="13"> {{.CompanyName}} of {{.CompanyAddress}}, hereinafter known as</text-chunk>
<text-chunk font="times" font-size="13"> the “Landlord”.</text-chunk>
</paragraph>

<paragraph margin="20 0 0 0">
<text-chunk font="times-bold" font-size="13">AND</text-chunk>
</paragraph>

<paragraph margin="20 0 0 0">
<text-chunk font="times" font-size="13">2 individuals known as Alex Tenant and Joanna Tenant, hereinafter known as the “Tenant(s)”, agree to the following:</text-chunk>
</paragraph>

<paragraph margin="20 0 0 0">
<text-chunk font="times-bold" font-size="13">OCCUPANT(S): </text-chunk>
<text-chunk font="times" font-size="13">The Premises is to be occupied strictly as a residential dwelling with the
following Two (2) Occupants to reside on the Premises in addition to the Tenant(s) mentioned above: Alex Jr Tenant and Jill Tenant, hereinafter known as the “Occupant(s)”.</text-chunk>
</paragraph>

<paragraph margin="20 0 0 0">
<text-chunk font="times-bold" font-size="13">OFFER TO RENT: </text-chunk>
<text-chunk font="times" font-size="13">The Landlord hereby rents to the Tenant(s), subject to the following terms and conditions of this Agreement, an apartment with the address of 1 Main Street, Apt 4, Small Town, Alabama, 20992 consisting of 2.5 bathroom(s) and 2 bedroom(s) hereinafter known as the “Premises”. The Landlord may also use the address for notices sent to the Tenant(s).</text-chunk>
</paragraph>

<paragraph margin="20 0 0 0">
<text-chunk font="times-bold" font-size="13">PURPOSE: </text-chunk>
<text-chunk font="times" font-size="13">The Tenant(s) and any Occupant(s) may only use the Premises as a residential
dwelling. It may not be used for storage, manufacturing of any type of food or product,
professional service(s), or for any commercial use unless otherwise stated in this Agreement.</text-chunk>
</paragraph>

<paragraph margin="20 0 0 0">
<text-chunk font="times-bold" font-size="13">FURNISHINGS: </text-chunk>
<text-chunk font="times" font-size="13">The Premises is furnished with the following:</text-chunk>
</paragraph>

<paragraph margin="20 0 0 0">
<text-chunk font="times" font-size="13">Bedroom Set(s), Dining Room Set(s), Living Room Set(s) and all other furnishings to be provided by the Tenant(s). Any damage to the Landlord's furnishings shall be the liability of the Tenant(s), reasonable wear-and-tear excepted, to be billed directly or less the Security Deposit.</text-chunk>
</paragraph>

<paragraph margin="20 0 0 0">
<text-chunk font="times-bold" font-size="13">APPLIANCES: </text-chunk>
<text-chunk font="times" font-size="13">The Landlord shall provide the following appliances:</text-chunk>
</paragraph>

<paragraph margin="20 0 0 0">
<text-chunk font="times" font-size="13">Air Conditioner(s), Dishwasher, Dryer (for Laundry), Fan(s), Hot Water Heater, HVAC,
Microwave, Outdoor Grill, Oven(s), Refrigerator, Stove(s), Washer (for Laundry), and any
other unnamed appliances existing on the Premises. Any damage to the Landlord's appliances
shall be the liability of the Tenant(s), reasonable wear-and-tear excepted, to be billed directly or less the Security Deposit.</text-chunk>
 </paragraph>
 
<paragraph margin="20 0 0 0">
<text-chunk font="times-bold" font-size="13">LEASE TERM: </text-chunk>
<text-chunk font="times" font-size="13"> This Agreement shall be a fixed-period arrangement beginning on December 03 2020 and ending on November 29 2033 with the Tenant(s) having the option to continue to occupy the Premises under the same terms and conditions of this Agreement under a Month-to-Month arrangement (Tenancy at Will) with either the Landlord or Tenant having the option to cancel the tenancy with at least thirty (30) days notice or the minimum time-period set by the State, whichever is shorter. For the Tenant to continue under Month-to-Month tenancy at the expiration of the Lease Term, the Landlord must be notified within sixty (60) days before the end of the Lease Term. Hereinafter known as the “Lease Term”.</text-chunk>
</paragraph>

<paragraph margin="20 0 0 0">
<text-chunk font="times-bold" font-size="13">RENT: </text-chunk>
<text-chunk font="times" font-size="13">Tenant(s) shall pay the Landlord in equal monthly installments of $1,873.00 (US
Dollars) hereinafter known as the “Rent”. The Rent will be due on the First (1st) of every
month and be paid through an electronic payment known as Automated Clearing House or “ACH”. Details of the Tenant's banking information and authorization shall be attached to this Lease Agreement.</text-chunk>
</paragraph>


<paragraph margin="20 0 0 0">
<text-chunk font="times-bold" font-size="13">NON-SUFFICIENT FUNDS (NSF CHECKS):</text-chunk>
<text-chunk font="times" font-size="13">If the Tenant(s) attempts to pay the rent with a check that is not honored or an electronic transaction (ACH) due to insufficient funds (NSF) there shall be a fee of $45.00 (US Dollars)</text-chunk>
</paragraph>
<paragraph margin="20 0 0 0">
<text-chunk font="times-bold" font-size="13">LATE FEE: </text-chunk>
<text-chunk font="times" font-size="13">If rent is not paid on the due date, there shall be a late fee assessed by the
Landlord in the amount of:

$50.00 (US Dollars) per occurrence for each month payment that is late after the 3rd Day rent
is due.
</text-chunk>
</paragraph>

<paragraph margin="20 0 0 0">
<text-chunk font="times-bold" font-size="13">FIRST (1ST) MONTH'S RENT:</text-chunk>
<text-chunk font="times" font-size="13">First (1st) month's rent shall be due by the Tenant(s) upon the execution of this Agreement.</text-chunk>
</paragraph>

<paragraph margin="20 0 0 0">
<text-chunk font="times-bold" font-size="13">PRE-PAYMENT:</text-chunk>
<text-chunk font="times" font-size="13">The Landlord shall not require any pre-payment of rent by the Tenant(s).</text-chunk>
</paragraph>

<paragraph margin="20 0 0 0">
<text-chunk font="times-bold" font-size="13">PROBATION PERIOD:</text-chunk>
<text-chunk font="times" font-size="13">The Tenant(s) will not move into the Premises before the start of the Lease Term.
</text-chunk>
</paragraph>

<paragraph margin="20 0 0 0">
<text-chunk font="times-bold" font-size="13">SECURITY DEPOSIT:</text-chunk>
<text-chunk font="times" font-size="13">A Security Deposit in the amount of $1,873.00 (US Dollars) shall be
required by the Tenant(s) at the execution of this Agreement to the Landlord for the faithful
performance of all the terms and conditions. The Security Deposit is to be returned to the
Tenant(s) within 14 days after this Agreement has terminated, less any damage charges and
without interest. This Security Deposit shall not be credited towards rent unless the Landlord
gives their written consent.
</text-chunk>
</paragraph>

<paragraph margin="20 0 0 0">
<text-chunk font="times-bold" font-size="13">POSSESSION:</text-chunk>
<text-chunk font="times" font-size="13">enant(s) has examined the condition of the Premises and by taking possession acknowledges that they have accepted the Premises in good order and in its current 
condition except as herein otherwise stated. Failure of the Landlord to deliver possession of the Premises at the start of the Lease Term to the Tenant(s) shall terminate this Agreement at the option of the Tenant(s). Furthermore, under such failure to deliver possession by the Landlord, and if the Tenant(s) cancels this Agreement, the Security Deposit (if any) shall be returned to the Tenant(s) along with any other pre-paid rent, fees, including if the Tenant(s) 
paid a fee during the application process before the execution of this Agreement.
</text-chunk>
</paragraph>


<paragraph margin="20 0 0 0">
<text-chunk font="times-bold" font-size="13">OPTION TO PURCHASE.</text-chunk>
<text-chunk font="times" font-size="13"> The Tenant(s) shall have the right to purchase the Premises
described herein for $450,000.00 at any time during the course of the Lease Term, along with
any renewal periods or extensions, by providing written notice to the Landlord along with a
deposit of $4,500.00 that is subject to the terms and conditions of a Purchase and Sale
Agreement to be negotiated, in “good faith”, between the Landlord and Tenant(s).
</text-chunk>
</paragraph>

<paragraph margin="20 0 0 0">
<text-chunk font="times" font-size="13">If the Landlord and Tenant(s) cannot produce a signed Purchase and Sale Agreement within a
reasonable time period then the deposit shall be refunded to the Tenant(s) and this Lease
Agreement shall continue under its terms and conditions.
</text-chunk>
</paragraph>

<paragraph margin="20 0 0 0">
<text-chunk font="times" font-size="13">If the option to purchase is exercised by the Tenant(s) all Rent that is paid to the Landlord shall remain separate from any and all deposits, consideration, or payments, made to the Landlord in regards to the purchase of the Premises.
</text-chunk>
</paragraph>

<paragraph margin="20 0 0 0">
<text-chunk font="times-bold" font-size="13">RECORDING. </text-chunk>
<text-chunk font="times" font-size="13">The Tenant(s) shall be withheld from recording this Option to Purchase unless the Tenant(s) has the written consent from the Landlord.
</text-chunk>
</paragraph>

<paragraph margin="20 0 0 0">
<text-chunk font="times-bold" font-size="13">ACCESS:</text-chunk>
<text-chunk font="times" font-size="13">TUpon the beginning of the Proration Period or the start of the Lease Term,
whichever is earlier, the Landlord agrees to give access to the Tenant(s) in the form of keys,
fobs, cards, or any type of keyless security entry as needed to enter the common areas and the
Premises. Duplicate copies of the access provided may only be authorized under the consent of the Landlord and, if any replacements are needed, the Landlord may provide them for a fee. At the end of this Agreement all access provided to the Tenant(s) shall be returned to the Landlord or a fee will be charged to the Tenant(s) or the fee will be subtracted from the Security Deposit.
</text-chunk>
</paragraph>

<paragraph margin="20 0 0 0">
<text-chunk font="times-bold" font-size="13">MOVE-IN INSPECTION:</text-chunk>
<text-chunk font="times" font-size="13">Before, at the time of the Tenant(s) accepting possession, or
shortly thereafter, the Landlord and Tenant(s) shall perform an inspection documenting the
present condition of all appliances, fixtures, furniture, and any existing damage within the
Premises.
</text-chunk>
</paragraph>


<paragraph margin="20 0 0 0">
<text-chunk font="times-bold" font-size="13">SUBLETTING:</text-chunk>
<text-chunk font="times" font-size="13">The Tenant(s) shall not have the right to sub-let the Premises or any part
thereof without the prior written consent of the Landlord. If consent is granted by the Landlord,
the Tenant(s) will be responsible for all actions and liabilities of the Sublessee including but not
limited to: damage to the Premises, non-payment of rent, and any eviction process (In the event
of an eviction the Tenant(s) shall be responsible for all court filing fee(s), representation, and
any other fee(s) associated with removing the Sublessee). The consent by the Landlord to one
sub-let shall not be deemed to be consent to any subsequent subletting.
</text-chunk>
</paragraph>


<paragraph margin="20 0 0 0">
<text-chunk font="times-bold" font-size="13">ABANDONMENT:</text-chunk>
<text-chunk font="times" font-size="13">If the Tenant(s) vacates or abandons the property for a time-period that is
the minimum set by State law or seven (7) days, whichever is less, the Landlord shall have the
right to terminate this Agreement immediately and remove all belongings including any
personal property off of the Premises. If the Tenant(s) vacates or abandons the property, the
Landlord shall immediately have the right to terminate this Agreement.
</text-chunk>
</paragraph>

<paragraph margin="20 0 0 0">
<text-chunk font="times-bold" font-size="13">ASSIGNMENT:</text-chunk>
<text-chunk font="times" font-size="13">Tenant(s) shall not assign this Lease without the prior written consent of the
Landlord. The consent by the Landlord to one assignment shall not be deemed to be consent to
any subsequent assignment.
</text-chunk>
</paragraph>

<paragraph margin="20 0 0 0">
<text-chunk font="times-bold" font-size="13">PARKING:</text-chunk>
<text-chunk font="times" font-size="13">The Landlord shall provide the Tenant(s) 2 Parking Spaces.
</text-chunk>
</paragraph>


<paragraph margin="20 0 0 0">
<text-chunk font="times" font-size="13">The Landlord shall not charge a fee for the 2 Parking Spaces. The Parking Space(s) can be
described as: 1 outdoor parking space and 1 indoor garage parking space provided
</text-chunk>
</paragraph>


<paragraph margin="20 0 0 0">
<text-chunk font="times-bold" font-size="13">RIGHT OF ENTRY:</text-chunk>
<text-chunk font="times" font-size="13">The Landlord shall have the right to enter the Premises during normal
working hours by providing notice in accordance with the minimum State requirement in order
for inspection, make necessary repairs, alterations or improvements, to supply services as
agreed or for any reasonable purpose. The Landlord may exhibit the Premises to prospective
purchasers, mortgagees, or lessees upon reasonable notice.
</text-chunk>
</paragraph>

<paragraph margin="20 0 0 0">
<text-chunk font="times-bold" font-size="13">SALE OF PROPERTY:</text-chunk>
<text-chunk font="times" font-size="13">If the Premises is sold, the Tenant(s) is to be notified of the new
Owner, and if there is a new Manager, their contact details for repairs and maintenance shall be
forwarded. If the Premises is conveyed to another party, the new owner shall not have the right
to terminate this Agreement and it shall continue under the terms and conditions agreed upon
by the Landlord and Tenant(s).
</text-chunk>
</paragraph>

<paragraph margin="20 0 0 0">
<text-chunk font="times-bold" font-size="13">UTILITIES:</text-chunk>
<text-chunk font="times" font-size="13">The Landlord agrees to pay for the following utilities and services:
</text-chunk>
</paragraph>

<paragraph margin="20 0 0 0">
<text-chunk font="times" font-size="13">
Lawn Care, Snow Removal, Trash Removal, Water, and the Landlord shall also provideSome
great services with all other utilities and services to be the responsibility of the Tenant(s).
</text-chunk>
</paragraph>

<paragraph margin="20 0 0 0">
<text-chunk font="times-bold" font-size="13">MAINTENANCE, REPAIRS, OR ALTERATIONS:</text-chunk>
<text-chunk font="times" font-size="13"> The Tenant(s) shall, at their own
expense and at all times, maintain the Premises in a clean and sanitary manner, and shall
surrender the same at termination hereof, in as good condition as received, normal wear and
tear excepted. The Tenant(s) may not make any alterations to the leased Premises without the
consent in writing of the Landlord. The Landlord shall be responsible for repairs to the interior
and exterior of the building. If the Premises includes a washer, dryer, freezer, dehumidifier unit
and/or air conditioning unit, the Landlord makes no warranty as to the repair or replacement of
units if one or all shall fail to operate. The Landlord will place fresh batteries in all
battery-operated smoke detectors when the Tenant(s) moves into the Premises. After the initial
placement of the fresh batteries, it is the responsibility of the Tenant(s) to replace batteries
when needed. A monthly “cursory” inspection may be required for all fire extinguishers to
make sure they are fully charged.
</text-chunk>
</paragraph>

<paragraph margin="20 0 0 0">
<text-chunk font="times-bold" font-size="13">EARLY TERMINATION:</text-chunk>
<text-chunk font="times" font-size="13"> The Tenant(s) may be allowed to cancel this Agreement under the
following conditions:
</text-chunk>
</paragraph>

<paragraph margin="20 0 0 0">
<text-chunk font="times" font-size="13">
The Tenant(s) must provide at least 60 days' notice and pay an early termination fee of
$1,000.00 (US Dollars) which does not include the rent due for the notice period. During the
notice period of 60 days the rent shall be paid in accordance with this Agreement.
</text-chunk>
</paragraph>

<paragraph margin="20 0 0 0">
<text-chunk font="times-bold" font-size="13">
PETS:
</text-chunk>
<text-chunk font="times" font-size="13">The Tenant(s) shall be allowed to have:
</text-chunk>
</paragraph>

<paragraph margin="10 0 0 0">
<text-chunk font="times" font-size="13">
Two (2) pets on the Premises consisting of Birds, Cats, Dogs, Fish, Hamsters, Rabbits, with no
other types of Pet(s) being allowed on the Premises or common areas, hereinafter known as the
“Pet(s)”. The Tenant(s) shall be required to pay a pet fee in the amount of $300.00 for all the
Pet(s) which is refundable at the end of the Lease Term only if there is no damage to the
Premises that is caused by the Pet(s). The Tenant(s) is responsible for all damage that any pet
causes, regardless of ownership of said pet and agrees to restore the property to its original
condition at their expense. There shall be no limit on the weight of the pet. pounds (Lb.).
</text-chunk>
</paragraph>

<paragraph margin="20 0 0 0">
<text-chunk font="times-bold" font-size="13">NOISE/WASTE: </text-chunk>
<text-chunk font="times" font-size="13">The Tenant(s) agrees not to commit waste on the Premises, maintain, or
permit to be maintained, a nuisance thereon, or use, or permit the Premises to be used, in an
unlawful manner. The Tenant(s) further agrees to abide by any and all local, county, and State
noise ordinances.
</text-chunk>
</paragraph>

<paragraph margin="20 0 0 0">
<text-chunk font="times-bold" font-size="13">GUESTS: </text-chunk>
<text-chunk font="times" font-size="13">There shall be no other persons living on the Premises other than the Tenant(s) and
any Occupant(s). Guests of the Tenant(s) are allowed for periods not lasting for more than
forty-eight hours unless otherwise approved by the Landlord.
</text-chunk>
</paragraph>


<paragraph margin="20 0 0 0">
<text-chunk font="times-bold" font-size="13">SMOKING POLICY: </text-chunk>
<text-chunk font="times" font-size="13">Smoking on the Premises is prohibited on the entire property, including
individual units, common areas, every building and adjoining properties.
MULTIPLE TENANT(S) OR OCCUPANT(S): Each individual that is considered a
Tenant(s) is jointly and individually liable for all of this Agreement's obligations, including but
not limited to rent monies. If any Tenant(s), guest, or Occupant(s) violates this Agreement, the
Tenant(s) is considered to have violated this Agreement. Landlord’s requests and notices to the
Tenant(s) or any of the Occupant(s) of legal age constitutes notice to the Tenant(s). Notices and
</text-chunk>
</paragraph>

<paragraph margin="20 0 0 0">
<text-chunk font="times-bold" font-size="13">COMPLIANCE WITH LAW: </text-chunk>
<text-chunk font="times" font-size="13">The Tenant(s) agrees that during the term of the Agreement, to
promptly comply with any present and future laws, ordinances, orders, rules, regulations, and
requirements of the Federal, State, County, City, and Municipal government or any of their
departments, bureaus, boards, commissions and officials thereof with respect to the Premises,
or the use or occupancy thereof, whether said compliance shall be ordered or directed to or
against the Tenant(s), the Landlord, or both.
</text-chunk>
</paragraph>

<paragraph margin="20 0 0 0">
<text-chunk font="times-bold" font-size="13">DEFAULT: </text-chunk>
<text-chunk font="times" font-size="13">If the Tenant(s) fails to comply with any of the financial or material provisions of
this Agreement, or of any present rules and regulations or any that may be hereafter prescribed
by the Landlord, or materially fails to comply with any duties imposed on the Tenant(s) by
statute or State laws, within the time period after delivery of written notice by the Landlord
specifying the non-compliance and indicating the intention of the Landlord to terminate the
Agreement by reason thereof, the Landlord may terminate this Agreement. If the Tenant(s) fails
to pay rent when due and the default continues for the time-period specified in the written
notice thereafter, the Landlord may, at their option, declare the entire balance (compiling all
months applicable to this Agreement) of rent payable hereunder to be immediately due and
payable and may exercise any and all rights and remedies available to the Landlord at law or in
equity and may immediately terminate this Agreement.
</text-chunk>
</paragraph>

<paragraph margin="20 0 0 0">
<text-chunk font="times" font-size="13">
The Tenant(s) will be in default if: (a) Tenant(s) does not pay rent or other amounts that are
owed in accordance with respective State laws; (b) Tenant(s), their guests, or the Occupant(s)
violate this Agreement, rules, or fire, safety, health, or criminal laws, regardless of whether
arrest or conviction occurs; (c) Tenant(s) abandons the Premises; (d) Tenant(s) gives incorrect
or false information in the rental application; (e) Tenant(s), or any Occupant(s) is arrested,
convicted, or given deferred adjudication for a criminal offense involving actual or potential
physical harm to a person, or involving possession, manufacture, or delivery of a controlled
substance, marijuana, or drug paraphernalia under state statute; (f) any illegal drugs or
paraphernalia are found in the Premises or on the person of the Tenant(s), guests, or
Occupant(s) while on the Premises and/or; (g) as otherwise allowed by law
</text-chunk>
</paragraph>

<paragraph margin="20 0 0 0">
<text-chunk font="times-bold" font-size="13">MULTIPLE TENANT(S) OR OCCUPANT(S): </text-chunk>
<text-chunk font="times" font-size="13">Each individual that is considered a Tenant(s) is jointly and individually liable for all of this Agreement's obligations, including but
not limited to rent monies. If any Tenant(s), guest, or Occupant(s) violates this Agreement, the
Tenant(s) is considered to have violated this Agreement. Landlord’s requests and notices to the
Tenant(s) or any of the Occupant(s) of legal age constitutes notice to the Tenant(s). Notices and requests from the Tenant(s) or any one of the Occupant(s) (including repair requests and entry
permissions) constitutes notice from the Tenant(s). In eviction suits, the Tenant(s) is considered
the agent of the Premise for the service of process.
</text-chunk>
</paragraph>

<paragraph margin="20 0 0 0">
<text-chunk font="times-bold" font-size="13">DISPUTES: </text-chunk>
<text-chunk font="times" font-size="13">If a dispute arises during or after the term of this Agreement between the
Landlord and Tenant(s), they shall agree to hold negotiations amongst themselves, in “good
faith”, before any litigation.
</text-chunk>
</paragraph>


<paragraph margin="20 0 0 0">
<text-chunk font="times-bold" font-size="13">SEVERABILITY: </text-chunk>
<text-chunk font="times" font-size="13">If any provision of this Agreement or the application thereof shall, for any
reason and to any extent, be invalid or unenforceable, neither the remainder of this Agreement
nor the application of the provision to other persons, entities or circumstances shall be affected
thereby, but instead shall be enforced to the maximum extent permitted by law.
</text-chunk>
</paragraph>

<paragraph margin="20 0 0 0">
<text-chunk font="times-bold" font-size="13">SURRENDER OF PREMISES: </text-chunk>
<text-chunk font="times" font-size="13">The Tenant(s) has surrendered the Premises when (a) the moveout date has passed and no one is living in the Premise within the Landlord’s reasonable
judgment; or (b) Access to the Premise have been turned in to Landlord – whichever comes
first. Upon the expiration of the term hereof, the Tenant(s) shall surrender the Premise in better
or equal condition as it were at the commencement of this Agreement, reasonable use, wear and
tear thereof, and damages by the elements excepted.
</text-chunk>
</paragraph>

<paragraph margin="20 0 0 0">
<text-chunk font="times-bold" font-size="13">RETALIATION: </text-chunk>
<text-chunk font="times" font-size="13">The Landlord is prohibited from making any type of retaliatory acts against
the Tenant(s) including but not limited to restricting access to the Premises, decreasing or
cancelling services or utilities, failure to repair appliances or fixtures, or any other type of act
that could be considered unjustified.
</text-chunk>
</paragraph>

<paragraph margin="20 0 0 0">
<text-chunk font="times-bold" font-size="13">WAIVER: </text-chunk>
<text-chunk font="times" font-size="13">A Waiver by the Landlord for a breach of any covenant or duty by the Tenant(s),
under this Agreement is not a waiver for a breach of any other covenant or duty by the
Tenant(s), or of any subsequent breach of the same covenant or duty. No provision of this
Agreement shall be considered waived unless such a waiver shall be expressed in writing as a
formal amendment to this Agreement and executed by the Tenant(s) and Landlord.
</text-chunk>
</paragraph>

<paragraph margin="20 0 0 0">
<text-chunk font="times-bold" font-size="13">EQUAL HOUSING: </text-chunk>
<text-chunk font="times" font-size="13">If the Tenant(s) possess(es) any mental or physical impairment, the
Landlord shall provide reasonable modifications to the Premises unless the modifications
would be too difficult or expensive for the Landlord to provide. Any impairment of the
Tenant(s) is/are encouraged to be provided and presented to the Landlord in writing in order to
seek the most appropriate route for providing the modifications to the Premises.
</text-chunk>
</paragraph>

<paragraph margin="20 0 0 0">
<text-chunk font="times-bold" font-size="13">HAZARDOUS MATERIALS: </text-chunk>
<text-chunk font="times" font-size="13">The Tenant(s) agrees to not possess any type of personal
property that could be considered a fire hazard such as a substance having flammable or explosive characteristics on the Premises. Items that are prohibited to be brought into the
Premises, other than for everyday cooking or the need of an appliance, includes but is not
limited to gas (compressed), gasoline, fuel, propane, kerosene, motor oil, fireworks, or any
other related content in the form of a liquid, solid, or gas.
</text-chunk>
</paragraph>

<paragraph margin="20 0 0 0">
<text-chunk font="times-bold" font-size="13">WATERBEDS: </text-chunk>
<text-chunk font="times" font-size="13">The Tenant(s) is not permitted to furnish the Premises with waterbeds.
</text-chunk>
</paragraph>

<paragraph margin="20 0 0 0">
<text-chunk font="times-bold" font-size="13">INDEMNIFICATION: </text-chunk>
<text-chunk font="times" font-size="13">The Landlord shall not be liable for any damage or injury to the
Tenant(s), or any other person, or to any property, occurring on the Premises, or any part
thereof, or in common areas thereof, and the Tenant(s) agrees to hold the Landlord harmless
from any claims or damages unless caused solely by the Landlord's negligence. It is
recommended that renter's insurance be purchased at the Tenant(s)'s expense.
</text-chunk>
</paragraph>

<paragraph margin="20 0 0 0">
<text-chunk font="times-bold" font-size="13">COVENANTS: </text-chunk>
<text-chunk font="times" font-size="13">The covenants and conditions herein contained shall apply to and bind the
heirs, legal representatives, and assigns of the parties hereto, and all covenants are to be
construed as conditions of this Agreement
</text-chunk>
</paragraph>

<paragraph margin="20 0 0 0">
<text-chunk font="times-bold" font-size="13">NOTICES: </text-chunk>
<text-chunk font="times" font-size="13">Any notice to be sent by the Landlord or the Tenant(s) to each other shall use the
following mailing addresses:
</text-chunk>
</paragraph>

<paragraph margin="20 0 0 0">
<text-chunk font="times-bold" font-size="13">Landlord's/Agent's Mailing Address</text-chunk>
</paragraph>

<paragraph margin="20 0 0 0">
<text-chunk font="times" font-size="13">
Best Landlord Company, ATTN. John Landlord
2 Maple Ln, Suite A, Best Town, Alabama, 29227
</text-chunk>
</paragraph>

<paragraph margin="20 0 0 0">
<text-chunk font="times-bold" font-size="13">
Tenant(s)'s Mailing Address
</text-chunk>
</paragraph>

<paragraph margin="20 0 0 0">
<text-chunk font="times" font-size="13">
Alex Tenant and Joanna Tenant
1 Main Street, Apt 4, Small Town, Alabama, 20992</text-chunk>
</paragraph>

<paragraph margin="20 0 0 0">
<text-chunk font="times-bold" font-size="13">AGENT/MANAGER: </text-chunk>
<text-chunk font="times" font-size="13">The Landlord authorizes the following to act on their behalf in regards
to the Premises for any repair, maintenance, or compliant other than a breach of this
Agreement: The The management company known as Best Management Company of 5 Maple
Ave, Suite 12A, Best City, Alabama, 29277 that can be contacted at the following Phone
Number (888) 222-3333 and can be E-Mailed at email@email.com.
</text-chunk>
</paragraph>

<paragraph margin="20 0 0 0">
<text-chunk font="times-bold" font-size="13">PREMISES DEEMED UNINHABITABLE: </text-chunk>
<text-chunk font="times" font-size="13">If the Property is deemed uninhabitable due to
damage beyond reasonable repair the Tenant(s) will be able to terminate this Agreement by
written notice to the Landlord. If said damage was due to the negligence of the Tenant(s), the
Tenant(s) shall be liable to the Landlord for all repairs and for the loss of income due to
restoring the Premises back to a livable condition in addition to any other losses that can be
proved by the Landlord.
</text-chunk>
</paragraph>

<paragraph margin="20 0 0 0">
<text-chunk font="times-bold" font-size="13">SERVICEMEMBERS CIVIL RELIEF ACT: </text-chunk>
<text-chunk font="times" font-size="13">In the event the Tenant(s) is or hereafter
becomes, a member of the United States Armed Forces on extended active duty and hereafter
the Tenant(s) receives permanent change of station (PCS) orders to depart from the area where
the Premises are located, or is relieved from active duty, retires or separates from the military,
is ordered into military housing, or receives deployment orders, then in any of these events, the
Tenant may terminate this lease upon giving thirty (30) days written notice to the Landlord.
The Tenant shall also provide to the Landlord a copy of the official orders or a letter signed by
the Tenant’s commanding officer, reflecting the change which warrants termination under this
clause. The Tenant will pay prorated rent for any days which he/she occupies the dwelling past
the beginning of the rental period.
</text-chunk>
<text-chunk font="times" font-size="13" margin="20 0 0 0">
The damage/security deposit will be promptly returned to Tenant, provided there are no
damages to the Premises
</text-chunk>
</paragraph>

<paragraph margin="20 0 0 0">
<text-chunk font="times-bold" font-size="13">LEAD PAINT: </text-chunk>
<text-chunk font="times" font-size="13">The Premises was not constructed before 1978 and therefore does not contain
leadbased paint.
</text-chunk>
</paragraph>

<paragraph margin="20 0 0 0">
<text-chunk font="times-bold" font-size="13">GOVERNING LAW: </text-chunk>
<text-chunk font="times" font-size="13">This Agreement is to be governed under the laws located in the State of
North Carolina
</text-chunk>
</paragraph>

<paragraph margin="20 0 0 0">
<text-chunk font="times-bold" font-size="13">ADDITIONAL TERMS AND CONDITIONS: </text-chunk>
<text-chunk font="times" font-size="13">In addition to the above stated terms and
conditions of this Agreement, the Landlord and Tenant agree to the following: Additional
Terms are to be specified: Term 1, Term 2, Term 3
</text-chunk>
</paragraph>

<paragraph margin="20 0 0 0">
<text-chunk font="times-bold" font-size="13">ENTIRE AGREEMENT: </text-chunk>
<text-chunk font="times" font-size="13">This Agreement contains all the terms agreed to by the parties
relating to its subject matter including any attachments or addendums. This Agreement replaces
all previous discussions, understandings, and oral agreements. The Landlord and Tenant(s)
agree to the terms and conditions and shall be bound until the end of the Lease Term.
</text-chunk>
<text-chunk font="times" font-size="13" margin="20 0 0 0">
The parties have agreed and executed this agreement on December 09 2020.
</text-chunk>
</paragraph>

<paragraph margin="20 0 10 0">
<text-chunk font="times-bold" font-size="13">
LANDLORD(S) SIGNATURE
</text-chunk>
</paragraph>

{{template "form-sig" dict "Margin" 110 "Text" "Landlord’s Signature"}}
<paragraph margin = "5 0 0 0">
<text-chunk font="times" font-size="13">John Landlord as President of Best Landlord Company</text-chunk>
</paragraph>

<paragraph margin="20 0 20 0">
<text-chunk font="times-bold" font-size="13">
TENANT(S) SIGNATURE
</text-chunk>
</paragraph>

{{template "form-sig" dict "Margin" "0 0 0 100" "Text" "Tenant’s Signature"}}

<division margin="20 0 20 0">
{{template "form-sig" dict "Margin" "0 0 0 100" "Text" "Tenant’s Signature"}}
</division>

<paragraph margin="0 0 0 0" text-align = "center">
   <text-chunk font="times-bold" font-size="21.5">Security Deposit Receipt</text-chunk>
</paragraph>

<division margin="10 0 0 0">
<paragraph line-height="2.5">
<text-chunk>
Date:__________________________________________  
Dear ___________________________________________[Tenant(s)],
The Landlord shall hold the Security Deposit in a separate account at a bank
located at __________________________________________[Street Address] in
the City of __________________________________________ , State of ___________________
The Security Deposit in the amount of $ _____________________ (US Dollars) has been deposited in
___________________ [Bank Name] with the Account Number of _________________ for the full
performance of the Lease executed on the _____ day of _______________ , 20 ___.
Sincerely,                                    .
</text-chunk>
<text-chunk font="times-bold" font-size="12"> Landlord’s Signature _________________________</text-chunk>
</paragraph>
</division>
<paragraph line-height="2.5">
</paragraph>
<division>
<paragraph margin="30 0 30 0" text-align = "center">
   <text-chunk font="times" font-size="13">AMOUNT ($) DUE AT SIGNING</text-chunk>
</paragraph>
</division>

<paragraph margin="30 0 30 0" text-align = "left" line-height="2.3">
<text-chunk font="times-bold" font-size="12">Security Deposit: </text-chunk> 
<text-chunk font="times" font-size="12">$1,873.00</text-chunk>
   <text-chunk font="times-bold" font-size="12">
First (1st) Month's Rent: </text-chunk><text-chunk font="times" font-size="12"> $1,873.00</text-chunk>
<text-chunk font="times-bold" font-size="12">
Pet Fee(s):</text-chunk> 
<text-chunk font="times" font-size="12"> $300.00 for all the Pet(s)</text-chunk>
</paragraph>