import { describe, it, expect, vi } from 'vitest';
import { fireEvent, screen } from '@testing-library/react';
import { render } from '../test/test-utils';
import { AddInventoryModal } from './AddInventoryModal';
import type { EquipmentItem } from '../equipmentTypes';

const equipmentItem: EquipmentItem = {
  id: 'equip-1',
  name: 'Shop Motor',
  category: 'motors',
  manufacturer: 'Acme',
  price: 21.5,
  currency: 'USD',
  seller: 'Test Seller',
  sellerId: 'seller-1',
  productUrl: 'https://example.com/motor',
  inStock: true,
  stockQty: 10,
  lastChecked: '2026-02-11T00:00:00Z',
};

describe('AddInventoryModal validation', () => {
  it('rejects quantity 0 while adding new gear', async () => {
    const onSubmit = vi.fn().mockResolvedValue(undefined);

    render(
      <AddInventoryModal
        isOpen
        onClose={vi.fn()}
        onSubmit={onSubmit}
        equipmentItem={equipmentItem}
      />,
    );

    const quantityInput = screen.getAllByRole('spinbutton')[0] as HTMLInputElement;
    fireEvent.change(quantityInput, { target: { value: '0' } });
    fireEvent.click(screen.getByRole('button', { name: 'Add to Inventory' }));

    expect(quantityInput).toHaveAttribute('min', '1');
    expect(onSubmit).not.toHaveBeenCalled();
  });

  it('validates purchase price before submit', async () => {
    const onSubmit = vi.fn().mockResolvedValue(undefined);

    render(
      <AddInventoryModal
        isOpen
        onClose={vi.fn()}
        onSubmit={onSubmit}
        equipmentItem={equipmentItem}
      />,
    );

    const purchasePriceInput = screen.getAllByRole('spinbutton')[1] as HTMLInputElement;
    fireEvent.change(purchasePriceInput, { target: { value: '-1' } });
    fireEvent.click(screen.getByRole('button', { name: 'Add to Inventory' }));

    expect(await screen.findByText('Enter a valid purchase price')).toBeInTheDocument();
    expect(onSubmit).not.toHaveBeenCalled();
  });
});
